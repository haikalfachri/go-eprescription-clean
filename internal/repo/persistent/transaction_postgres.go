package persistent

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/pkg/postgres"
)

// TransactionRepo -.
type TransactionRepo struct {
	*postgres.Postgres
}

// NewTransactionRepo - creates a new Transaction repository.
func NewTransactionRepo(pg *postgres.Postgres) *TransactionRepo {
	return &TransactionRepo{pg}
}

// Create - inserts a new transaction record.
func (r *TransactionRepo) Create(ctx context.Context, t entity.Transaction) (*entity.Transaction, error) {
	sql, args, err := r.Builder.
		Insert("transactions").
		Columns("patient_id", "medicine_type", "status", "total_price", "total_medicines", "payment_token", "payment_redirect_url").
		Values(t.PatientID, t.MedicineType, "pending", t.TotalPrice, t.TotalMedicines, t.PaymentToken, t.PaymentRedirectURL).
		Suffix("RETURNING id, patient_id, medicine_type, status, total_price, total_medicines, payment_token, payment_redirect_url").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo - Create - Builder: %w", err)
	}

	var transaction entity.Transaction
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&transaction.ID, &transaction.PatientID, &transaction.MedicineType, &transaction.Status, &transaction.TotalPrice, &transaction.TotalMedicines, &transaction.PaymentToken, &transaction.PaymentRedirectURL)
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo - Create - QueryRow: %w", err)
	}

	return &transaction, nil
}

// GetAll - retrieves all transaction records.
func (r *TransactionRepo) GetAll(ctx context.Context) ([]entity.Transaction, error) {
	sql, _, err := r.Builder.
		Select("id", "patient_id", "medicine_type", "status", "total_price", "total_medicines").
		From("transactions").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo - GetAll - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo - GetAll - Query: %w", err)
	}
	defer rows.Close()

	var transactions []entity.Transaction
	for rows.Next() {
		var transaction entity.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.PatientID, &transaction.MedicineType, &transaction.Status, &transaction.TotalPrice, &transaction.TotalMedicines); err != nil {
			return nil, fmt.Errorf("TransactionRepo - GetAll - Scan: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("TransactionRepo - GetAll - no transactions found")
	}

	return transactions, nil
}

// GetAllByPatientID - retrieves all transactions by patient ID.
func (r *TransactionRepo) GetAllByPatientID(ctx context.Context, patientID string) ([]entity.Transaction, error) {
	sql, args, err := r.Builder.
		Select("id", "patient_id", "medicine_type", "status", "total_price", "total_medicines", "payment_token", "payment_redirect_url").
		From("transactions").
		Where("patient_id = ?", patientID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo - GetAllByPatientID - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo - GetAllByPatientID - Query: %w", err)
	}
	defer rows.Close()

	var transactions []entity.Transaction
	for rows.Next() {
		var transaction entity.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.PatientID, &transaction.MedicineType, &transaction.Status, &transaction.TotalPrice, &transaction.TotalMedicines, &transaction.PaymentToken, &transaction.PaymentRedirectURL); err != nil {
			return nil, fmt.Errorf("TransactionRepo - GetAllByPatientID - Scan: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("TransactionRepo - GetAllByPatientID - no transactions found for patient ID %s", patientID)
	}

	return transactions, nil
}

// GetByID - retrieves a transaction by its ID.
func (r *TransactionRepo) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	sql, args, err := r.Builder.
		Select("id", "patient_id", "medicine_type", "status", "total_price", "total_medicines", "payment_token", "payment_redirect_url").
		From("transactions").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo - GetByID - Builder: %w", err)
	}

	var transaction entity.Transaction
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&transaction.ID, &transaction.PatientID, &transaction.MedicineType, &transaction.Status, &transaction.TotalPrice, &transaction.TotalMedicines, &transaction.PaymentToken, &transaction.PaymentRedirectURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("TransactionRepo - GetByID - no transaction found with ID %s", id)
		}
		return nil, fmt.Errorf("TransactionRepo - GetByID - QueryRow: %w", err)
	}

	return &transaction, nil
}

// Update - updates an existing transaction record.
func (r *TransactionRepo) Update(ctx context.Context, id string, t entity.Transaction) (*entity.Transaction, error) {
	updates := map[string]any{}

	if t.PatientID != "" {
		updates["patient_id"] = t.PatientID
	}
	if t.MedicineType != "" {
		updates["medicine_type"] = t.MedicineType
	}
	if t.TotalPrice != 0 {
		updates["total_price"] = t.TotalPrice
	}
	if t.TotalMedicines != 0 {
		updates["total_medicines"] = t.TotalMedicines
	}
	if t.Status != "" {
		updates["status"] = t.Status
	}
	if t.PaymentToken != "" {
		updates["payment_token"] = t.PaymentToken
	}
	if t.PaymentRedirectURL != "" {
		updates["payment_redirect_url"] = t.PaymentRedirectURL
	}

	// If no fields to update
	if len(updates) == 0 {
		return nil, fmt.Errorf("TransactionRepo - Update - no fields to update")
	}

	sql, args, err := r.Builder.
		Update("transactions").
		SetMap(updates).
		Where("id = ?", id).
		Suffix("RETURNING id, patient_id, medicine_type, status, total_price, total_medicines, payment_token, payment_redirect_url").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo - Update - Builder: %w", err)
	}

	var transaction entity.Transaction
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&transaction.ID, &transaction.PatientID, &transaction.MedicineType, &transaction.Status, &transaction.TotalPrice, &transaction.TotalMedicines , &transaction.PaymentToken, &transaction.PaymentRedirectURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("TransactionRepo - Update - no transaction found with ID %s", id)
		}
		return nil, fmt.Errorf("TransactionRepo - Update - QueryRow: %w", err)
	}

	return &transaction, nil
}

// Delete - deletes a transaction record by its ID.
func (r *TransactionRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.Builder.
		Delete("transactions").
		Where("id = ?").
		ToSql()
	if err != nil {
		return fmt.Errorf("TransactionRepo - Delete - Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, append(args, id)...)
	if err != nil {
		return fmt.Errorf("TransactionRepo - Delete - Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("TransactionRepo - Delete - no transaction found with ID %s", id)
	}

	return nil
}

// UpdateStatusByTransactionID updates the status of a transaction based on its ID (used as TransactionID in Midtrans)
func (r *TransactionRepo) UpdateStatusByTransactionID(ctx context.Context, transactionID, status string) error {
	sql, args, err := r.Builder.
		Update("transactions").
		Set("status", status).
		Where("id = ?", transactionID). // TransactionID maps to your transaction.ID
		ToSql()
	if err != nil {
		return fmt.Errorf("TransactionRepo - UpdateStatusByTransactionID - Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TransactionRepo - UpdateStatusByTransactionID - Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("TransactionRepo - UpdateStatusByTransactionID - no rows affected for ID %s", transactionID)
	}

	return nil
}
