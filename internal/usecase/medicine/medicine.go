package medicine

import (
	"context"
	"fmt"
	"strings"
	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/internal/repo"
)

// UseCase - Medicine use case struct.
type UseCase struct {
	repo repo.MedicineRepo
}

// NewUseCase - creates a new Medicine use case.
func New(repo repo.MedicineRepo) *UseCase {
	return &UseCase{repo}
}

// Create - stores a new medicine.
func (uc *UseCase) Create(ctx context.Context, m entity.Medicine) (*entity.Medicine, error) {
	medicine, err := uc.repo.Create(ctx, m)
	if err != nil {
		return nil, fmt.Errorf("MedicineUseCase - Create - repo.Create: %w", err)
	}
	return medicine, nil
}

// GetAll - retrieves all medicines.
func (uc *UseCase) GetAll(ctx context.Context) ([]entity.Medicine, error) {
	medicines, err := uc.repo.GetAll(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no medicines found") {
			return nil, fmt.Errorf("MedicineUseCase - GetAll - repo.GetAll: no medicines found")
		}
		return nil, fmt.Errorf("MedicineUseCase - GetAll - repo.GetAll: %w", err)
	}
	return medicines, nil
}

// GetByID - retrieves a medicine by ID.
func (uc *UseCase) GetByID(ctx context.Context, id string) (*entity.Medicine, error) {
	medicine, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "medicine not found") {
			return nil, fmt.Errorf("MedicineUseCase - GetByID - repo.GetByID: medicine not found")
		}
		return nil, fmt.Errorf("MedicineUseCase - GetByID - repo.GetByID: %w", err)
	}
	return medicine, nil	
}

// Update - updates an existing medicine.
func (uc *UseCase) Update(ctx context.Context, id string, m entity.Medicine	) (*entity.Medicine, error) {
	updatedMedicine, err := uc.repo.Update(ctx, id, m)
	if err != nil {
		if strings.Contains(err.Error(), "medicine not found") {
			return nil, fmt.Errorf("MedicineUseCase - Update - repo.Update: medicine not found")
		}
		return nil, fmt.Errorf("MedicineUseCase - Update - repo.Update: %w", err)
	}
	return updatedMedicine, nil
}

// Delete - deletes a medicine by ID.
func (uc *UseCase) Delete(ctx context.Context, id string) error {	
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "medicine not found") {
			return fmt.Errorf("MedicineUseCase - Delete - repo.Delete: medicine not found")
		}
		return fmt.Errorf("MedicineUseCase - Delete - repo.Delete: %w", err)
	}
	return nil
}