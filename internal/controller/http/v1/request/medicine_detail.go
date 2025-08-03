package request

// MedicineDetail - represents the details of a medicine in a transaction.
type MedicineDetail struct {
	MedicineID	string   `json:"medicine_id" example:"123e4567-e89b-12d3-a456-426614174000"`  
	SignaID     string   `json:"signa_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Description string   `json:"description" example:"For pain relief"`
	Quantity    int64    `json:"quantity" validate:"required,gte=0" example:"2"`
}

type MedicineDetailID struct {
	ID string `validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`
}