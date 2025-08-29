package v1

import (
	v1 "go-eprescription-clean/docs/proto/v1"
	"go-eprescription-clean/pkg/logger"
	"github.com/go-playground/validator/v10"
	pbgrpc "google.golang.org/grpc"
	"go-eprescription-clean/pkg/rabbitmq/server"
)

// NewAuditRoutes -.
func NewAuditRoutes(app *pbgrpc.Server, rmq map[string]server.CallHandler, u Usecases, l logger.Interface) {
	r := &V1{u: u, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	{
		// gRPC service
		v1.RegisterAuditServiceServer(app, r)

		// RabbitMQ handler
		rmq["transaction.event"] = r.handleEvent()
	}
}