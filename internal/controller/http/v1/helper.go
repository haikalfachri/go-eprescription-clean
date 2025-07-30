package v1

import (
	"fmt"
	"strings"
	"go-eprescription-clean/internal/controller/http/v1/response"
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
)

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(response.Error{Error: msg})
}

func successResponse[T any](ctx *fiber.Ctx, code int, msg string, data T) error {
	return ctx.Status(code).JSON(response.Success[T]{
		Message: msg,
		Data:    data,
	})
}

func formatValidationErrors(ve validator.ValidationErrors) string {
	var messages []string

	for _, fe := range ve {
		field := fe.Field()
		tag := fe.Tag()

		var msg string
		switch tag {
		case "required":
			msg = fmt.Sprintf("%s is required", field)
		case "gte":
			msg = fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
		case "lte":
			msg = fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
		case "oneof":
			msg = fmt.Sprintf("%s must be one of [%s]", field, fe.Param())
		case "uuid4":
			msg = fmt.Sprintf("%s must be a valid uuid", field)
		default:
			msg = fmt.Sprintf("%s is not valid", field)
		}

		messages = append(messages, msg)
	}

	return strings.Join(messages, ", ")
}