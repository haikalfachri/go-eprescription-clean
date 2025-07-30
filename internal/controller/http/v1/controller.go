package v1

import (
	"github.com/go-playground/validator/v10"
	"go-eprescription-clean/internal/usecase"
	"go-eprescription-clean/pkg/logger"
)

// V1 -.
type Usecases struct {
	Signa       usecase.Signa
	Patient     usecase.Patient
	Medicine    usecase.Medicine
	Transaction usecase.Transaction
	MedicineDetail usecase.MedicineDetail
	// ... 
}

type V1 struct {
	u Usecases
	l logger.Interface
	v *validator.Validate
}
