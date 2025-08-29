package audit

import (
	"context"
	"fmt"
	"encoding/json"

	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/internal/repo"
)

// UseCase - Transaction use case struct.
type UseCase struct {
	repo repo.AuditRepo
}

// New - creates a new Transaction use case.
func New(r repo.AuditRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (u *UseCase) StoreEvent(ctx context.Context, body []byte) error {
	var req map[string]interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		return err
	}

	var auditLog entity.AuditLog

	if event, ok := req["event"].(string); ok {
		auditLog.Event = event
	}

	auditLog.Payload = string(body)

	_, err := u.repo.Create(ctx, auditLog)
	return err
}

func (u *UseCase) GetAll(ctx context.Context) ([]entity.AuditLog, error) {
	events, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}
	return events, nil
}


