package transaction

import (
	"context"
	"fmt"
	"strings"

	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/internal/repo"
)

// UseCase - Transaction use case struct.
type UseCase struct {
	repo               repo.TransactionRepo
	medicineDetailRepo repo.MedicineDetailRepo
	medicineRepo       repo.MedicineRepo
}

// New - creates a new Transaction use case.
func New(r repo.TransactionRepo, mdRepo repo.MedicineDetailRepo, mRepo repo.MedicineRepo) *UseCase {
	return &UseCase{
		repo:               r,
		medicineDetailRepo: mdRepo,
		medicineRepo:       mRepo,
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
		return nil, fmt.Errorf("TransactionUseCase - CreateWithMedicineDetail - repo.CreateWithMedicineDetail: the length of medicines, signas, quantities, and descriptions  must be equal")
	}
	// 1. Create the transaction
	transaction, err := uc.repo.Create(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("TransactionUseCase - CreateWithMedicineDetail - repo.CreateWithMedicineDetail: failed to create transaction: %w", err)
	}

	// 2. Insert medicine details with transaction ID
	for idx := range medicines {
		var medicineDetail entity.MedicineDetail
		medicineDetail.TransactionID = transaction.ID
		medicineDetail.MedicineID = medicines[idx]
		medicineDetail.SignaID = signas[idx]
		medicineDetail.Description = descriptions[idx]
		medicineDetail.Quantity = quantities[idx]
		medicine, err := uc.medicineRepo.GetByID(ctx, medicines[idx])
		if err != nil {
			return nil, fmt.Errorf("TransactionUseCase - CreateWithMedicineDetail - repo.GetByID: failed to get medicine by ID: %w", err)
		}
		transaction.TotalPrice += quantities[idx] * medicine.Price
		medicine.Quantity -= quantities[idx]
		if _, err := uc.medicineRepo.Update(ctx, medicines[idx], *medicine); err != nil {
			return nil, fmt.Errorf("TransactionUseCase - CreateWithMedicineDetail - repo.Update: failed to update medicine quantity: %w", err)
		}
		if _, err := uc.medicineDetailRepo.Create(ctx, medicineDetail); err != nil {
			return nil, fmt.Errorf("TransactionUseCase - CreateWithMedicineDetail - repo.CreateWithMedicineDetail: failed to create medicine detail: %w", err)
		}
	}

	transaction.TotalMedicines = int64(len(medicines))

	transaction, err = uc.repo.Update(ctx, transaction.ID, *transaction)
	if err != nil {
		return nil, fmt.Errorf("TransactionUseCase - CreateWithMedicineDetail - repo.Update: failed to update transaction: %w", err)
	}

	details, err := uc.medicineDetailRepo.GetByTransactionID(ctx, transaction.ID)
	if err != nil {
		return nil, fmt.Errorf("TransactionUseCase - CreateWithMedicineDetail - GetByTransactionID: %w", err)
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
