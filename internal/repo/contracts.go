// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"go-eprescription-clean/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/test/mocks_repo_test.go -package=usecase_test

type (
	// SignaRepo -.
	SignaRepo interface {
		Create(context.Context, entity.Signa) (*entity.Signa, error)
		GetAll(context.Context) ([]entity.Signa, error)
		GetByID(context.Context, string) (*entity.Signa, error)
		Update(context.Context, string, entity.Signa) (*entity.Signa, error)
		Delete(context.Context, string) error
	}

	// PatientRepo -.
	PatientRepo interface {
		Create(context.Context, entity.Patient) (*entity.Patient, error)
		GetAll(context.Context) ([]entity.Patient, error)
		GetByID(context.Context, string) (*entity.Patient, error)
		Update(context.Context, string, entity.Patient) (*entity.Patient, error)
		Delete(context.Context, string) error
	}

	// MedicineRepo -.
	MedicineRepo interface {
		Create(context.Context, entity.Medicine) (*entity.Medicine, error)
		GetAll(context.Context) ([]entity.Medicine, error)
		GetByID(context.Context, string) (*entity.Medicine, error)
		Update(context.Context, string, entity.Medicine) (*entity.Medicine, error)
		Delete(context.Context, string) error
	}

	// TransactionRepo -.
	TransactionRepo interface {
		Create(context.Context, entity.Transaction) (*entity.Transaction, error)
		GetAll(context.Context) ([]entity.Transaction, error)
		GetAllByPatientID(context.Context, string) ([]entity.Transaction, error)
		GetByID(context.Context, string) (*entity.Transaction, error)
		Update(context.Context, string, entity.Transaction) (*entity.Transaction, error)
		Delete(context.Context, string) error
	}

	// MedicineDetailRepo -.
	MedicineDetailRepo interface {
		Create(context.Context, entity.MedicineDetail) (*entity.MedicineDetail, error)
		GetAll(context.Context) ([]entity.MedicineDetail, error)
		GetByTransactionID(context.Context, string) ([]entity.MedicineDetail, error)
		GetByID(context.Context, string) (*entity.MedicineDetail, error)
		Update(context.Context, string, entity.MedicineDetail) (*entity.MedicineDetail, error)
		Delete(context.Context, string) error
	}
)

