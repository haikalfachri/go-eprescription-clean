package v1

import (
	v1 "go-eprescription-clean/docs/proto/v1"
	"go-eprescription-clean/internal/usecase"
	"go-eprescription-clean/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// Usecases -.
type Usecases struct {
	Audit usecase.Audit
	// ...
}

// V1 -.
type V1 struct {
	v1.AuditServiceServer

	u Usecases
	l logger.Interface
	v *validator.Validate
}