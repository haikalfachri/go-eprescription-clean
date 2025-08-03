// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// MedicineDetail -.
type MedicineDetail struct {
	ID            string `json:"id" example:"1789dfb3-e060-433a-9db4-cb93b42768d8"`
	TransactionID string `json:"transaction_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	SignaID       string `json:"signa_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	MedicineID    string `json:"medicine_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Quantity	  int64  `json:"quantity" example:"2"`
	Description   string `json:"description" example:"For pain relief"`
}
