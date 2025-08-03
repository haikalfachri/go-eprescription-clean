package persistent

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/pkg/postgres"
)

// MedicinePostgresRepo - Postgres implementation of Medicine repository.
type MedicinePostgresRepo struct {
	*postgres.Postgres
}

// NewMedicineRepo - creates a new Medicine repository.
func NewMedicineRepo(pg *postgres.Postgres) *MedicinePostgresRepo {
	return &MedicinePostgresRepo{pg}
}

// Create - inserts a new medicine record.
func (r *MedicinePostgresRepo) Create(ctx context.Context, m entity.Medicine) (*entity.Medicine, error) {
	sql, args, err := r.Builder.
		Insert("mastermedicines").
		Columns("name", "quantity", "price").
		Values(m.Name, m.Quantity, m.Price).
		Suffix("RETURNING id, name, qty, price").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MedicineRepo - Create - Builder: %w", err)
	}

	var medicine entity.Medicine
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&medicine.ID, &medicine.Name, &medicine.Quantity, &medicine.Price)
	if err != nil {
		return nil, fmt.Errorf("MedicineRepo - Create - QueryRow: %w", err)
	}

	return &medicine, nil
}

// GetAll - retrieves all medicine records.
func (r *MedicinePostgresRepo) GetAll(ctx context.Context) ([]entity.Medicine, error) {
	sql, _, err := r.Builder.
		Select("id", "name", "quantity", "price").
		From("mastermedicines").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MedicineRepo - GetAll - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("MedicineRepo - GetAll - Query: %w", err)
	}
	defer rows.Close()

	var medicines []entity.Medicine
	for rows.Next() {
		var medicine entity.Medicine
		if err := rows.Scan(&medicine.ID, &medicine.Name, &medicine.Quantity, &medicine.Price); err != nil {
			return nil, fmt.Errorf("MedicineRepo - GetAll - Scan: %w", err)
		}
		medicines = append(medicines, medicine)
	}

	if len(medicines) == 0 {
		return nil, fmt.Errorf("MedicineRepo - GetAll - no medicines found")
	}

	return medicines, nil
}

// GetByID - retrieves a medicine by ID.
func (r *MedicinePostgresRepo) GetByID(ctx context.Context, id string) (*entity.Medicine, error) {
	sql, args, err := r.Builder.
		Select("id", "name", "quantity", "price").
		From("mastermedicines").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MedicineRepo - GetByID - Builder: %w", err)
	}

	var medicine entity.Medicine
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&medicine.ID, &medicine.Name, &medicine.Quantity, &medicine.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("MedicineRepo - GetByID - QueryRow: medicine not found")
		}
		return nil, fmt.Errorf("MedicineRepo - GetByID - QueryRow: %w", err)
	}

	return &medicine, nil
}

// Update - updates an existing medicine record.
func (r *MedicinePostgresRepo) Update(ctx context.Context, id string, m entity.Medicine) (*entity.Medicine, error) {
	sql, args, err := r.Builder.
		Update("mastermedicines").
		SetMap(map[string]any{
			"name":     m.Name,
			"quantity": m.Quantity,
			"price":    m.Price,
		}).Where("id = ?", id).
		Suffix("RETURNING id, name, quantity, price").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MedicineRepo - Update - Builder: %w", err)
	}

	var medicine entity.Medicine
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&medicine.ID, &medicine.Name, &medicine.Quantity, &medicine.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("MedicineRepo - Update - QueryRow: medicine not found")
		}
		return nil, fmt.Errorf("MedicineRepo - Update - QueryRow: %w", err)
	}

	return &medicine, nil
}

// Delete - deletes a medicine record by ID.
func (r *MedicinePostgresRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.Builder.
		Delete("mastermedicines").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("MedicineRepo - Delete - Builder: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, append(args, id)...)
	if err != nil {
		return fmt.Errorf("MedicineRepo - Delete - Exec: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("MedicineRepo - Delete - Exec: medicine not found")
	}

	return nil
}
