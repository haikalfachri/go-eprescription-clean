package persistent

import (
	"context"
	"fmt"
	"errors"

	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

// SignaRepo -.
type SignaRepo struct {
	*postgres.Postgres
}

// NewSignaRepo - creates a new Signa repository.
func NewSignaRepo(pg *postgres.Postgres) *SignaRepo {
	return &SignaRepo{pg}
}

// Create - inserts a new signa record.
func (r *SignaRepo) Create(ctx context.Context, s entity.Signa) (*entity.Signa, error) {
	sql, args, err := r.Builder.
		Insert("mastersignas").
		Columns("signa", "description").
		Values(s.Signa, s.Description).
		Suffix("RETURNING id, signa, description").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("SignaRepo - Create - Builder: %w", err)
	}

	var signa entity.Signa
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&signa.ID, &signa.Signa, &signa.Description)
	if err != nil {
		return nil, fmt.Errorf("SignaRepo - Create - QueryRow: %w", err)
	}

	return &signa, nil
}

// GetAll - retrieves all signa records.
func (r *SignaRepo) GetAll(ctx context.Context) ([]entity.Signa, error) {
	sql, _, err := r.Builder.
		Select("id", "signa", "description").
		From("mastersignas").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("SignaRepo - GetAll - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("SignaRepo - GetAll - Query: %w", err)
	}
	defer rows.Close()

	signas := make([]entity.Signa, 0)

	for rows.Next() {
		var s entity.Signa
		err = rows.Scan(&s.ID, &s.Signa, &s.Description)
		if err != nil {
			return nil, fmt.Errorf("SignaRepo - GetAll - Scan: %w", err)
		}
		signas = append(signas, s)
	}

	if len(signas) == 0 {
		return nil, fmt.Errorf("SignaRepo - GetAll - no signas found")
	}

	return signas, nil
}

// GetByID - retrieves a signa by its ID.
func (r *SignaRepo) GetByID(ctx context.Context, id string) (*entity.Signa, error) {
	sql, args, err := r.Builder.
		Select("id", "signa", "description").
		From("mastersignas").
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("SignaRepo - GetByID - Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	var s entity.Signa
	err = row.Scan(&s.ID, &s.Signa, &s.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("SignaRepo - GetByID - QueryRow: signa not found")
		}
		return nil, fmt.Errorf("SignaRepo - GetByID - QueryRow: %w", err)
	}

	return &s, nil
}

// Update - updates an existing signa record.
func (r *SignaRepo) Update(ctx context.Context, id string, s entity.Signa) (*entity.Signa, error) {
	sql, args, err := r.Builder.
		Update("mastersignas").
		SetMap(map[string]any{
			"signa":       s.Signa,
			"description": s.Description,
		}).
		Where("id = ?", id).
		Suffix("RETURNING id, signa, description").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("SignaRepo - Update - Builder: %w", err)
	}

	var updatedSigna entity.Signa
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&updatedSigna.ID, &updatedSigna.Signa, &updatedSigna.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("SignaRepo - Update - QueryRow: signa not found")
		}
		return nil, fmt.Errorf("SignaRepo - Update - QueryRow: %w", err)
	}

	return &updatedSigna, nil
}

// Delete - deletes a signa by its ID.
func (r *SignaRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.Builder.
		Delete("mastersignas").
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return fmt.Errorf("SignaRepo - Delete - Builder: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SignaRepo - Delete - Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("SignaRepo - Delete - Exec: signa not found")
	}

	return nil
}
