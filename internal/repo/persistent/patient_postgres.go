package persistent

import (
	"context"
	"errors"
	"fmt"

	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/pkg/postgres"

	"github.com/jackc/pgx/v5"
)

// PatientRepo -.
type PatientRepo struct {
	*postgres.Postgres
}

// NewPatientRepo - creates a new Patient repository.
func NewPatientRepo(pg *postgres.Postgres) *PatientRepo {
	return &PatientRepo{pg}
}

// Create - inserts a new patient record.
func (r *PatientRepo) Create(ctx context.Context, p entity.Patient) (*entity.Patient, error) {
	sql, args, err := r.Builder.
		Insert("patients").
		Columns("name", "age", "gender").
		Values(p.Name, p.Age, p.Gender).
		Suffix("RETURNING id, name, age, gender").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PatientRepo - Create - Builder: %w", err)
	}

	var patient entity.Patient
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Gender)
	if err != nil {
		return nil, fmt.Errorf("PatientRepo - Create - QueryRow: %w", err)
	}

	return &patient, nil
}

// GetAll - retrieves all patient records.
func (r *PatientRepo) GetAll(ctx context.Context) ([]entity.Patient, error) {
	sql, _, err := r.Builder.
		Select("id", "name", "age", "gender").
		From("patients").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PatientRepo - GetAll - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("PatientRepo - GetAll - Query: %w", err)
	}
	defer rows.Close()

	var patients []entity.Patient
	for rows.Next() {
		var patient entity.Patient
		if err := rows.Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Gender); err != nil {
			return nil, fmt.Errorf("PatientRepo - GetAll - Scan: %w", err)
		}
		patients = append(patients, patient)
	}

	if len(patients) == 0 {
		return nil, fmt.Errorf("PatientRepo - GetAll - no patients found")
	}

	return patients, nil
}

// GetByID - retrieves a patient record by ID.
func (r *PatientRepo) GetByID(ctx context.Context, id string) (*entity.Patient, error) {
	sql, args, err := r.Builder.
		Select("id", "name", "age", "gender").
		From("patients").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PatientRepo - GetByID - Builder: %w", err)
	}

	var patient entity.Patient
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Gender)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("PatientRepo - GetByID - QueryRow: patient not found")
		}
		return nil, fmt.Errorf("PatientRepo - GetByID - QueryRow: %w", err)
	}

	return &patient, nil
}

// Update - updates an existing patient record.
func (r *PatientRepo) Update(ctx context.Context, id string, p entity.Patient) (*entity.Patient, error) {
	sql, args, err := r.Builder.
		Update("patients").
		SetMap(map[string]any{
			"name":   p.Name,
			"age":    p.Age,
			"gender": p.Gender,
		}).
		Where("id = ?", id).
		Suffix("RETURNING id, name, age, gender").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PatientRepo - Update - Builder: %w", err)
	}

	var patient entity.Patient
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Gender)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("PatientRepo - GetByID - QueryRow: patient not found")
		}
		return nil, fmt.Errorf("PatientRepo - Update - QueryRow: %w", err)
	}

	return &patient, nil
}

// Delete - deletes a patient record by ID.
func (r *PatientRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.Builder.
		Delete("patients").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("PatientRepo - Delete - Builder: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PatientRepo - Delete - Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("PatientRepo - Delete - Exec: patient not found")
	}

	return nil
}
