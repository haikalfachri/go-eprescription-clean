package request

type Patient struct {
	Name   string `json:"name" validate:"required" example:"John Doe"`
	Age    int    `json:"age" validate:"gte=0,lte=120" example:"30"`
	Gender string `json:"gender" validate:"required,oneof=male female" example:"male"`
}

type PatientID struct {
	ID string `validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`
}
