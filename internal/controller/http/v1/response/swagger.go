package response

import "go-eprescription-clean/internal/entity"

// SuccessString is a simple string response wrapper (e.g. for delete)
type SuccessString struct {
	Message string `json:"message" example:"success"`
	Data    string `json:"data" example:"deleted"`
}

// SuccessSigna is a non-generic response wrapper for Signa entity
type SuccessSigna struct {
	Message string       `json:"message" example:"success"`
	Data    entity.Signa `json:"data"`
}

// SuccessSignaList is for list responses
type SuccessSignaList struct {
	Message string         `json:"message" example:"success"`
	Data    []entity.Signa `json:"data"`
}

// SuccessPatient is a non-generic response wrapper for Patient entity
type SuccessPatient struct {
	Message string         `json:"message" example:"success"`
	Data    entity.Patient `json:"data"`
}

// SuccessPatientList is for list responses
type SuccessPatientList struct {
	Message string           `json:"message" example:"success"`
	Data    []entity.Patient `json:"data"`
}

// SuccessMedicine is a non-generic response wrapper for Medicine entity
type SuccessMedicine struct {
	Message string          `json:"message" example:"success"`
	Data    entity.Medicine `json:"data"`
}

// SuccessMedicineList is for list responses
type SuccessMedicineList struct {
	Message string            `json:"message" example:"success"`
	Data    []entity.Medicine `json:"data"`
}

// SuccessTransaction is a non-generic response wrapper for Transaction entity
type SuccessTransaction struct {
	Message string           `json:"message" example:"success"`
	Data    entity.Transaction `json:"data"`
}

// SuccessTransactionList is for list responses
type SuccessTransactionList struct {
	Message string             `json:"message" example:"success"`
	Data    []entity.Transaction `json:"data"`
}

// SuccessMedicineDetail is a non-generic response wrapper for MedicineDetail entity
type SuccessMedicineDetail struct {
	Message string              `json:"message" example:"success"`
	Data    entity.MedicineDetail `json:"data"`
}

// SuccessMedicineDetailList is for list responses
type SuccessMedicineDetailList struct {
	Message string               `json:"message" example:"success"`
	Data    []entity.MedicineDetail `json:"data"`
}