package patient

import (
	"context"
	"fmt"
	"strings"

	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/internal/repo"
)

// UseCase - Patient use case struct.
type UseCase struct {
	repo repo.PatientRepo
}

// NewUseCase - creates a new Patient use case.
func New(repo repo.PatientRepo) *UseCase {
	return &UseCase{repo}
}

// Create - stores a new patient.
func (uc *UseCase) Create(ctx context.Context, p entity.Patient) (*entity.Patient, error) {
	patient, err := uc.repo.Create(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("PatientUseCase - Create - repo.Create: %w", err)
	}
	return patient, nil
}

// GetAll - retrieves all patients.
func (uc *UseCase) GetAll(ctx context.Context) ([]entity.Patient, error) {
	patients, err := uc.repo.GetAll(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no patients found") {
			return nil, fmt.Errorf("PatientUseCase - GetAll - repo.GetAll: no patients found")
		}
		return nil, fmt.Errorf("PatientUseCase - GetAll - repo.GetAll: %w", err)
	}
	return patients, nil
}

// GetByID - retrieves a patient by ID.
func (uc *UseCase) GetByID(ctx context.Context, id string) (*entity.Patient, error) {
	patient, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "patient not found") {
			return nil, fmt.Errorf("PatientUseCase - GetByID - repo.GetByID: patient not found")
		}
		return nil, fmt.Errorf("PatientUseCase - GetByID - repo.GetByID: %w", err)
	}
	return patient, nil	
}

// Update - updates an existing patient.
func (uc *UseCase) Update(ctx context.Context, id string, p entity.Patient) (*entity.Patient, error) {
	updatedPatient, err := uc.repo.Update(ctx, id, p)
	if err != nil {
		if strings.Contains(err.Error(), "patient not found") {
			return nil, fmt.Errorf("PatientUseCase - Update - repo.Update: patient not found")
		}
		return nil, fmt.Errorf("PatientUseCase - Update - repo.Update: %w", err)
	}
	return updatedPatient, nil
}

// Delete - deletes a patient by ID.
func (uc *UseCase) Delete(ctx context.Context, id string) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "patient not found") {
			return fmt.Errorf("PatientUseCase - Delete - repo.Delete: patient not found")
		}
		return fmt.Errorf("PatientUseCase - Delete - repo.Delete: %w", err)
	}
	return nil
}

