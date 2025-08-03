package transaction

import (
	"context"
	"fmt"
	"strings"

	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/internal/repo"
	"go-eprescription-clean/pkg/midtrans"
)

// UseCase - Transaction use case struct.
type UseCase struct {
	repo               repo.TransactionRepo
	medicineDetailRepo repo.MedicineDetailRepo
	medicineRepo       repo.MedicineRepo
	patientRepo        repo.PatientRepo
	signaRepo          repo.SignaRepo
	midtransClient     *midtrans.SnapClient
}

// New - creates a new Transaction use case.
func New(r repo.TransactionRepo, mdRepo repo.MedicineDetailRepo, mRepo repo.MedicineRepo, pRepo repo.PatientRepo, sRepo repo.SignaRepo, mtClient *midtrans.SnapClient) *UseCase {
	return &UseCase{
		repo:               r,
		medicineDetailRepo: mdRepo,
		medicineRepo:       mRepo,
		patientRepo:        pRepo,
		signaRepo:          sRepo,
		midtransClient:     mtClient,
	}
}

// Create with MedicineDetail - stores a new transaction with medicine details.
func (uc *UseCase) CreateWithMedicineDetail(
	ctx context.Context,
	t entity.Transaction,
	medicines []string,
	signas []string,
	descriptions []string,
	quantities []int64,
) (*entity.Transaction, error) {
	if len(medicines) != len(signas) || len(signas) != len(descriptions) || len(quantities) != len(medicines) {
		return nil, fmt.Errorf("all input arrays must have the same length")
	}

	// 1. Validate patient
	patient, err := uc.patientRepo.GetByID(ctx, t.PatientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find patient: %w", err)
	}

	// 2. Validate all medicines and signas first
	var medicineDetails []entity.MedicineDetail
	totalPrice := int64(0)

	for idx := range medicines {
		medicine, err := uc.medicineRepo.GetByID(ctx, medicines[idx])
		if err != nil {
			return nil, fmt.Errorf("medicine ID %s not found: %w", medicines[idx], err)
		}

		if medicine.Quantity < quantities[idx] {
			return nil, fmt.Errorf("insufficient stock for medicine ID %s", medicines[idx])
		}

		_, err = uc.signaRepo.GetByID(ctx, signas[idx])
		if err != nil {
			return nil, fmt.Errorf("signa ID %s not found: %w", signas[idx], err)
		}

		// prepare medicine detail
		md := entity.MedicineDetail{
			MedicineID:  medicines[idx],
			SignaID:     signas[idx],
			Description: descriptions[idx],
			Quantity:    quantities[idx],
		}
		medicineDetails = append(medicineDetails, md)
		totalPrice += quantities[idx] * medicine.Price
	}

	// 3. Create transaction (safe to create now)
	t.TotalPrice = totalPrice
	t.TotalMedicines = int64(len(medicines))
	transaction, err := uc.repo.Create(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 4. Create medicine details and update stock
	for _, detail := range medicineDetails {
		detail.TransactionID = transaction.ID

		medicine, _ := uc.medicineRepo.GetByID(ctx, detail.MedicineID)
		medicine.Quantity -= detail.Quantity

		if _, err := uc.medicineRepo.Update(ctx, detail.MedicineID, *medicine); err != nil {
			return nil, fmt.Errorf("failed to update medicine quantity: %w", err)
		}

		if _, err := uc.medicineDetailRepo.Create(ctx, detail); err != nil {
			return nil, fmt.Errorf("failed to create medicine detail: %w", err)
		}
	}

	// 5. Create snap transaction
	snapResp, _ := uc.midtransClient.CreateSnapTransaction(transaction.ID, totalPrice, patient.Name)
	if snapResp == nil {
		return nil, fmt.Errorf("failed to create snap transaction")
	}

	transaction.PaymentToken = snapResp.Token
	transaction.PaymentRedirectURL = snapResp.RedirectURL

	// 6. Save payment token
	transaction, err = uc.repo.Update(ctx, transaction.ID, *transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction with snap info: %w", err)
	}

	// 7. Fetch details and return
	details, err := uc.medicineDetailRepo.GetByTransactionID(ctx, transaction.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get medicine details: %w", err)
	}
	transaction.MedicineDetail = details

	return transaction, nil
}

// GetAll - returns all transactions.
func (uc *UseCase) GetAll(ctx context.Context) ([]entity.Transaction, error) {
	transactions, err := uc.repo.GetAll(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no transactions found") {
			return nil, fmt.Errorf("TransactionUseCase - GetAll - repo.GetAll: no transactions found")
		}
		return nil, fmt.Errorf("TransactionUseCase - GetAll - repo.GetAll: %w", err)
	}

	transactions, err = uc.attachMedicineDetails(ctx, transactions)
	if err != nil {
		return nil, fmt.Errorf("TransactionUseCase - GetAll - attachMedicineDetails: %w", err)
	}
	return transactions, nil
}

// GetAllByPatientID - returns all transactions by patient ID.
func (uc *UseCase) GetAllByPatientID(ctx context.Context, patientID string) ([]entity.Transaction, error) {
	transactions, err := uc.repo.GetAllByPatientID(ctx, patientID)
	if err != nil {
		if strings.Contains(err.Error(), "no transactions found") {
			return nil, fmt.Errorf("TransactionUseCase - GetAllByPatientID - repo.GetAllByPatientID: no transactions found")
		}
		return nil, fmt.Errorf("TransactionUseCase - GetAllByPatientID - repo.GetAllByPatientID: %w", err)
	}

	transactions, err = uc.attachMedicineDetails(ctx, transactions)
	if err != nil {
		return nil, fmt.Errorf("TransactionUseCase - GetAllByPatientID - attachMedicineDetails: %w", err)
	}
	return transactions, nil
}

// GetByID - returns a transaction by ID.
func (uc *UseCase) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	transaction, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "transaction not found") {
			return nil, fmt.Errorf("TransactionUseCase - GetByID - repo.GetByID: transaction not found")
		}
		return nil, fmt.Errorf("TransactionUseCase - GetByID - repo.GetByID: %w", err)
	}

	details, err := uc.medicineDetailRepo.GetByTransactionID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("TransactionUseCase - GetByID - GetByTransactionID: %w", err)
	}
	transaction.MedicineDetail = details

	return transaction, nil
}

// Update - updates a transaction by ID.
func (uc *UseCase) Update(ctx context.Context, id string, t entity.Transaction) (*entity.Transaction, error) {
	updatedTransaction, err := uc.repo.Update(ctx, id, t)
	if err != nil {
		if strings.Contains(err.Error(), "transaction not found") {
			return nil, fmt.Errorf("TransactionUseCase - Update - repo.Update: transaction not found")
		}
		return nil, fmt.Errorf("TransactionUseCase - Update - repo.Update: %w", err)
	}

	details, err := uc.medicineDetailRepo.GetByTransactionID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("TransactionUseCase - Update - GetByTransactionID: %w", err)
	}
	updatedTransaction.MedicineDetail = details

	return updatedTransaction, nil
}

// Delete - deletes a transaction by ID.
func (uc *UseCase) Delete(ctx context.Context, id string) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "transaction not found") {
			return fmt.Errorf("TransactionUseCase - Delete - repo.Delete: transaction not found")
		}
		return fmt.Errorf("TransactionUseCase - Delete - repo.Delete: %w", err)
	}
	return nil
}

// attachMedicineDetails - attaches medicine details to transactions.
func (uc *UseCase) attachMedicineDetails(ctx context.Context, transactions []entity.Transaction) ([]entity.Transaction, error) {
	for idx, transaction := range transactions {
		details, err := uc.medicineDetailRepo.GetByTransactionID(ctx, transaction.ID)
		if err != nil {
			return nil, fmt.Errorf("TransactionUseCase - attachMedicineDetails: %w", err)
		}
		transactions[idx].MedicineDetail = details
	}
	return transactions, nil
}

// HandleMidtransNotification - handles Midtrans notification callback.
func (uc *UseCase) HandleMidtransNotification(ctx context.Context, transactionID, transactionStatus, fraudStatus string) error {
	// You can map status to your internal status system
	var internalStatus string
	switch transactionStatus {
	case "capture":
		if fraudStatus == "challenge" {
			internalStatus = "pending"
		} else if fraudStatus == "accept" {
			internalStatus = "paid"
		}
	case "settlement":
		internalStatus = "paid"
	case "deny", "cancel", "expire":
		internalStatus = "failed"
	case "pending":
		internalStatus = "pending"
	}

	// Update the transaction status in your database
	return uc.repo.UpdateStatusByTransactionID(ctx, transactionID, internalStatus)
}
