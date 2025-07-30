// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"go-eprescription-clean/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=./test/mocks_usecase_test.go -package=usecase_test
type (

	// Signa -.
	Signa interface {
		Create(context.Context, entity.Signa) (*entity.Signa, error)
		GetAll(context.Context) ([]entity.Signa, error)
		GetByID(context.Context, string) (*entity.Signa, error)
		Update(context.Context, string, entity.Signa) (*entity.Signa, error)
		Delete(context.Context, string) (error)
	}

	// Patient -.
	Patient interface {
		Create(context.Context, entity.Patient) (*entity.Patient, error)
		GetAll(context.Context) ([]entity.Patient, error)
		GetByID(context.Context, string) (*entity.Patient, error)
		Update(context.Context, string, entity.Patient) (*entity.Patient, error)
		Delete(context.Context, string) error
	}	

	// Medicine -.
	Medicine interface {
		Create(context.Context, entity.Medicine) (*entity.Medicine, error)
		GetAll(context.Context) ([]entity.Medicine, error)
		GetByID(context.Context, string) (*entity.Medicine, error)
		Update(context.Context, string, entity.Medicine) (*entity.Medicine, error)
		Delete(context.Context, string) error
	}

	// Transaction -.
	Transaction interface {
		CreateWithMedicineDetail(context.Context, entity.Transaction, []string, []string, []string, []int64) (*entity.Transaction, error)
		GetAll(context.Context) ([]entity.Transaction, error)
		GetAllByPatientID(context.Context, string) ([]entity.Transaction, error)
		GetByID(context.Context, string) (*entity.Transaction, error)
		Update(context.Context, string, entity.Transaction) (*entity.Transaction, error)
		Delete(context.Context, string) error
	}

	// MedicineDetail -.
	MedicineDetail interface {
		Create(context.Context, entity.MedicineDetail) (*entity.MedicineDetail, error)
		GetAll(context.Context) ([]entity.MedicineDetail, error)
		GetByTransactionID(context.Context, string) ([]entity.MedicineDetail, error)
		GetByID(context.Context, string) (*entity.MedicineDetail, error)
		Update(context.Context, string, entity.MedicineDetail) (*entity.MedicineDetail, error)
		Delete(context.Context, string) error	
	}
)
