// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// Transaction-.
type Transaction struct {
	ID                 string           `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PatientID          string           `json:"patient_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	MedicineType       string           `json:"medicine_type" example:"compound"`
	TotalPrice         int64            `json:"total_price" example:"15000"`
	TotalMedicines     int64            `json:"total_medicines" example:"3"`
	Status             string           `json:"status" example:"pending"`
	MedicineDetail     []MedicineDetail `json:"medicine_detail,omitempty"`
	PaymentRedirectURL string           `json:"payment_redirect_url,omitempty" example:"https://example.com/redirect"`
	PaymentToken       string           `json:"payment_token,omitempty" example:"https://example.com/token"`
}
