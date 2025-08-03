package persistent

import (
	"context"
	"fmt"

	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

// MedicineDetailsRepo - repository for medicine details.
type MedicineDetailsRepo struct {
	*postgres.Postgres
}

// NewMedicineDetailsRepo - creates a new MedicineDetails repository.
func NewMedicineDetailsRepo(pg *postgres.Postgres) *MedicineDetailsRepo {
	return &MedicineDetailsRepo{pg}
}

// Create - inserts a new medicine details record.
func (r *MedicineDetailsRepo) Create(ctx context.Context, md entity.MedicineDetail) (*entity.MedicineDetail, error) {
	sql, args, err := r.Builder.
		Insert("medicinedetails").
		Columns("medicine_id", "transaction_id", "signa_id", "quantity", "description").
		Values(md.MedicineID, md.TransactionID, md.SignaID, md.Quantity, md.Description).
		Suffix("RETURNING id, medicine_id, transaction_id, signa_id, quantity, description").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MedicineDetailsRepo - Create - Builder: %w", err)
	}

	var medicineDetail entity.MedicineDetail
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&medicineDetail.ID, &medicineDetail.MedicineID, &medicineDetail.TransactionID, &medicineDetail.SignaID, &medicineDetail.Quantity, &medicineDetail.Description)
	if err != nil {
		return nil, fmt.Errorf("MedicineDetailsRepo - Create - QueryRow: %w", err)
	}

	return &medicineDetail, nil
}

// GetAll - retrieves all medicine details records.
func (r *MedicineDetailsRepo) GetAll(ctx context.Context) ([]entity.MedicineDetail, error) {
	sql, _, err := r.Builder.
		Select("id", "medicine_id", "transaction_id", "signa_id", "quantity", "description").
		From("medicinedetails").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MedicineDetailsRepo - GetAll - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("MedicineDetailsRepo - GetAll - Query: %w", err)
	}
	defer rows.Close()

	var medicineDetails []entity.MedicineDetail
	for rows.Next() {
		var md entity.MedicineDetail
		if err := rows.Scan(&md.ID, &md.MedicineID, &md.TransactionID, &md.SignaID, &md.Quantity, &md.Description); err != nil {
			return nil, fmt.Errorf("MedicineDetailsRepo - GetAll - Scan: %w", err)
		}
		medicineDetails = append(medicineDetails, md)
	}

	if len(medicineDetails) == 0 {
		return nil, fmt.Errorf("MedicineDetailsRepo - GetAll - no medicine details found")
	}

	return medicineDetails, nil
}

// GetByTransactionID - retrieves medicine details by transaction ID.
func (r *MedicineDetailsRepo) GetByTransactionID(ctx context.Context, transactionID string) ([]entity.MedicineDetail, error) {
	sql, args, err := r.Builder.
		Select("id", "medicine_id", "transaction_id", "signa_id", "description", "quantity").
		From("medicinedetails").
		Where("transaction_id = ?", transactionID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MedicineDetailsRepo - GetByTransactionID - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("MedicineDetailsRepo - GetByTransactionID - Query: %w", err)
	}
	defer rows.Close()

	var medicineDetails []entity.MedicineDetail
	for rows.Next() {
		var md entity.MedicineDetail
		if err := rows.Scan(&md.ID, &md.MedicineID, &md.TransactionID, &md.SignaID, &md.Description, &md.Quantity); err != nil {
			return nil, fmt.Errorf("MedicineDetailsRepo - GetByTransactionID - Scan: %w", err)
		}
		medicineDetails = append(medicineDetails, md)
	}

	if len(medicineDetails) == 0 {
		return nil, fmt.Errorf("MedicineDetailsRepo - GetByTransactionID - no medicine details found for transaction ID %s", transactionID)
	}

	return medicineDetails, nil
}

// GetByID - retrieves a medicine detail by its ID.
func (r *MedicineDetailsRepo) GetByID(ctx context.Context, id string) (*entity.MedicineDetail, error) {
	sql, args, err := r.Builder.
		Select("id", "medicine_id", "transaction_id", "signa_id", "description", "quantity").
		From("medicinedetails").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MedicineDetailsRepo - GetByID - Builder: %w", err)
	}

	var medicineDetail entity.MedicineDetail
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&medicineDetail.ID, &medicineDetail.MedicineID, &medicineDetail.TransactionID, &medicineDetail.SignaID, &medicineDetail.Description, &medicineDetail.Quantity)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("MedicineDetailsRepo - GetByID - no medicine detail found with ID %s", id)
		}
		return nil, fmt.Errorf("MedicineDetailsRepo - GetByID - QueryRow: %w", err)
	}

	return &medicineDetail, nil
}

// Update - updates an existing medicine detail record.
func (r *MedicineDetailsRepo) Update(ctx context.Context, id string, md entity.MedicineDetail) (*entity.MedicineDetail, error) {
	sql, args, err := r.Builder.
		Update("medicinedetails").
		Set("id_medicine", md.MedicineID).
		Set("id_transaction", md.TransactionID).
		Set("signa_id", md.SignaID).
		Set("description", md.Description).
		Where("id = ?", id).
		Suffix("RETURNING id, id_medicine, id_transaction, signa_id, description, quantity").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MedicineDetailsRepo - Update - Builder: %w", err)
	}

	var updatedMD entity.MedicineDetail
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&updatedMD.ID, &updatedMD.MedicineID, &updatedMD.TransactionID, &updatedMD.SignaID, &updatedMD.Description, &updatedMD.Quantity)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("MedicineDetailsRepo - Update - no medicine detail found with ID %s", md.ID)
		}
		return nil, fmt.Errorf("MedicineDetailsRepo - Update - QueryRow: %w", err)
	}

	return &updatedMD, nil
}

// Delete - removes a medicine detail by its ID.
func (r *MedicineDetailsRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.Builder.
		Delete("medicinedetails").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("MedicineDetailsRepo - Delete - Builder: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("MedicineDetailsRepo - Delete - Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("MedicineDetailsRepo - Delete - no medicine detail found with ID %s", id)
	}

	return nil
}
