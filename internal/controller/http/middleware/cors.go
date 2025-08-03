package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go-eprescription-clean/pkg/logger"
)

func CORS(l logger.Interface) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*", // Or specify: "https://example.com"
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS, PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	})
}
