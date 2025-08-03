package request

// Medicine -.
type Medicine struct {
	Name     string `json:"name" validate:"required" example:"Panadol"`
	Quantity int64  `json:"quantity" validate:"required,gte=0" example:"100"`
	Price    int64  `json:"price"    validate:"required,gte=0" example:"10000"`
}

type MedicineID struct {
	ID string `validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`
}

