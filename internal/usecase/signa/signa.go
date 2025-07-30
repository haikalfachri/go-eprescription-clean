package signa

import (
	"context"
	"fmt"
	"strings"

	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/internal/repo"
)

// UseCase - Signa use case struct.
type UseCase struct {
	repo repo.SignaRepo
}

// New - creates a new Signa use case.
func New(r repo.SignaRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

// Create - stores a new signa.
func (uc *UseCase) Create(ctx context.Context, s entity.Signa) (*entity.Signa, error) {
	signa, err := uc.repo.Create(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("SignaUseCase - Create - repo.Create: %w", err)
	}
	return signa, nil
}


// GetAll - returns all signas.
func (uc *UseCase) GetAll(ctx context.Context) ([]entity.Signa, error) {
	signas, err := uc.repo.GetAll(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no signas found") {
			return nil, fmt.Errorf("SignaUseCase - GetAll - repo.GetAll: no signas found")
		}
		return nil, fmt.Errorf("SignaUseCase - GetAll - repo.GetAll: %w", err)
	}
	return signas, nil
}

// GetByID - returns signa by id.
func (uc *UseCase) GetByID(ctx context.Context, id string) (*entity.Signa, error) {
	signa, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "signa not found") {
			return nil, fmt.Errorf("SignaUseCase - GetByID - repo.GetByID: signa not found")
		}
		return nil, fmt.Errorf("SignaUseCase - GetByID - repo.GetByID: %w", err)
	}
	return signa, nil
}

// Update - updates a signa by id.
func (uc *UseCase) Update(ctx context.Context, id string, s entity.Signa) (*entity.Signa, error) {
	updatedSigna, err := uc.repo.Update(ctx, id, s)
	if err != nil {
		if strings.Contains(err.Error(), "signa not found") {
			return nil, fmt.Errorf("SignaUseCase - Update - repo.Update: signa not found")
		}
		return nil, fmt.Errorf("SignaUseCase - Update - repo.Update: %w", err)
	}
	return updatedSigna, nil
}


// Delete - deletes a signa by id.
func (uc *UseCase) Delete(ctx context.Context, id string) error {
	err := uc.repo.Delete(ctx, id); 
	if err != nil {
		if strings.Contains(err.Error(), "signa not found") {
			return fmt.Errorf("SignaUseCase - Delete - repo.Delete: signa not found")
		}
		return fmt.Errorf("SignaUseCase - Delete - repo.Delete: %w", err)
	}
	return nil
}

