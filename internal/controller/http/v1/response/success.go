package response

type Success[T any] struct {
	Message string `json:"message" example:"operation successful"`
	Data    T      `json:"data"`
}
