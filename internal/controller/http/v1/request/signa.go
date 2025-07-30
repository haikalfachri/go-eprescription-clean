package request

type Signa struct {
	Signa       string `json:"signa"        example:"a.c"        validate:"required"`
	Description string `json:"description"  example:"Before eat" validate:"required"`
}

type SignaID struct {
	ID string `json:"id" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`
}
