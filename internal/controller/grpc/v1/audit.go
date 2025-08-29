package v1

import (
	"context"
	"time"

	v1 "go-eprescription-clean/docs/proto/v1"
	"go-eprescription-clean/pkg/rabbitmq/server"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Handle incoming audit log events from RabbitMQ .-
func (r *V1) handleEvent() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		r.l.Info("Received message: %s", d.Body)

		// Just pass raw body to usecase
		if err := r.u.Audit.StoreEvent(context.Background(), d.Body); err != nil {
			r.l.Error("failed to store event:", err)
			return nil, err
		}

		return string(d.Body), nil
	}
}

// Retrieve all audit logs .-
// Name of function must be the same as defined in protobuf file
func (r *V1) GetAllAuditLogs(ctx context.Context, _ *v1.GetLogRequest) (*v1.GetLogResponse, error) {
	logs, err := r.u.Audit.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var protoLogs []*v1.AuditLog
	for _, log := range logs {
		protoLogs = append(protoLogs, &v1.AuditLog{
			Id:            log.ID,
			Event:         log.Event,
			Payload:       log.Payload.(string),
			Timestamp:     log.Timestamp.Format(time.RFC3339),
		})
	}

	return &v1.GetLogResponse{Logs: protoLogs}, nil
}
