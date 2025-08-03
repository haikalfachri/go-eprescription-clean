package request

type CreateTransaction struct {
	PatientID    string   `json:"patient_id" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`
	MedicineType string   `json:"medicine_type" validate:"required,oneof=compound non_compound" example:"compound"`
	Medicines    []string `json:"medicines" validate:"required,dive,required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`
	Signas       []string `json:"signas" validate:"required,dive,required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`
	Quantities   []int64  `json:"quantities" validate:"required,dive,required,gte=0" example:"2"`
	Descriptions []string `json:"descriptions" validate:"required,dive,required" example:"For pain relief"`
}

type UpdateTransaction struct {
	MedicineType   string `json:"medicine_type" validate:"required,oneof=compound non-compound" example:"compound"`
	TotalPrice     int64  `json:"total_price" validate:"required,gte=0" example:"15000"`
	TotalMedicines int64  `json:"total_medicines" validate:"required,gte=0" example:"3"`
}

type TransactionID struct {
	ID string `json:"id" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`
}
