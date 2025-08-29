package message_broker

import (
	"context"
	"fmt"

	"go-eprescription-clean/pkg/rabbitmq/client"
)

type RMQRepo struct {
	client *client.Client
}

// NewRMQRepo initializes RMQRepo with your custom rabbitmq.Connection
func NewRMQRepo(c *client.Client) *RMQRepo {
	return &RMQRepo{client: c}
}

// PublishEvent publishes a message to the configured exchange
func (r *RMQRepo) PublishEvent(ctx context.Context, handler string, payload interface{}) error {
	var response interface{}
	err := r.client.RemoteCall(handler, payload, &response)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}
	return nil
}

