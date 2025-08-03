package medicine_detail

import (
	"context"
	"fmt"
	"strings"

	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/internal/repo"
)

// UseCase - MedicineDetail use case struct.
type UseCase struct {
	repo repo.MedicineDetailRepo
}

// New - creates a new MedicineDetail use case.
func New(r repo.MedicineDetailRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

// Create - stores a new medicine detail.
func (uc *UseCase) Create(ctx context.Context, md entity.MedicineDetail) (*entity.MedicineDetail, error) {
	medicineDetail, err := uc.repo.Create(ctx, md)
	if err != nil {
		return nil, fmt.Errorf("MedicineDetailUseCase - Create - repo.Create: %w", err)
	}
	return medicineDetail, nil
}

// GetAll - returns all medicine details.
func (uc *UseCase) GetAll(ctx context.Context) ([]entity.MedicineDetail, error) {
	medicineDetails, err := uc.repo.GetAll(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no medicine details found") {
			return nil, fmt.Errorf("MedicineDetailUseCase - GetAll - repo.GetAll: no medicine details found")
		}
		return nil, fmt.Errorf("MedicineDetailUseCase - GetAll - repo.GetAll: %w", err)
	}
	return medicineDetails, nil
}

// GetByTransactionID - retrieves medicine details by transaction ID.
func (uc *UseCase) GetByTransactionID(ctx context.Context, transactionID string) ([]entity.MedicineDetail, error) {
	medicineDetails, err := uc.repo.GetByTransactionID(ctx, transactionID)
	if err != nil {
		if strings.Contains(err.Error(), "no medicine details found for transaction ID") {
			return nil, fmt.Errorf("MedicineDetailUseCase - GetByTransactionID - repo.GetByTransactionID: %w", err)
		}
		return nil, fmt.Errorf("MedicineDetailUseCase - GetByTransactionID - repo.GetByTransactionID: %w", err)
	}
	return medicineDetails, nil
}

// GetByID - retrieves a medicine detail by its ID.
func (uc *UseCase) GetByID(ctx context.Context, id string) (*entity.MedicineDetail, error) {
	medicineDetail, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no medicine detail found for ID") {
			return nil, fmt.Errorf("MedicineDetailUseCase - GetByID - repo.GetByID: %w", err)
		}
		return nil, fmt.Errorf("MedicineDetailUseCase - GetByID - repo.GetByID: %w", err)
	}
	return medicineDetail, nil
}

// Update - updates a medicine detail by its ID.
func (uc *UseCase) Update(ctx context.Context, id string, md entity.MedicineDetail) (*entity.MedicineDetail, error) {
	updatedMedicineDetail, err := uc.repo.Update(ctx, id, md)
	if err != nil {
		if strings.Contains(err.Error(), "no medicine detail found for ID") {
			return nil, fmt.Errorf("MedicineDetailUseCase - Update - repo.Update: %w", err)
		}
		return nil, fmt.Errorf("MedicineDetailUseCase - Update - repo.Update: %w", err)
	}
	return updatedMedicineDetail, nil
}

// Delete - deletes a medicine detail by its ID.
func (uc *UseCase) Delete(ctx context.Context, id string) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no medicine detail found for ID") {
			return fmt.Errorf("MedicineDetailUseCase - Delete - repo.Delete: %w", err)
		}
		return fmt.Errorf("MedicineDetailUseCase - Delete - repo.Delete: %w", err)
	}
	return nil
}