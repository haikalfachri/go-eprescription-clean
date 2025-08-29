package grpc

import (
	v1 "go-eprescription-clean/internal/controller/grpc/v1"
	"go-eprescription-clean/pkg/logger"
	"go-eprescription-clean/pkg/rabbitmq/server"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewRouter -.
func NewRouter(app *pbgrpc.Server, u v1.Usecases, l logger.Interface)  map[string]server.CallHandler {
	rmqRoutes := make(map[string]server.CallHandler)
	{
		v1.NewAuditRoutes(app, rmqRoutes, u, l)
	}

	reflection.Register(app)

	return rmqRoutes
}
