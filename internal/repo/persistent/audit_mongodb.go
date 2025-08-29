package persistent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditRepo -.
type AuditRepo struct {
	*mongodb.MongoDB
	collection *mongo.Collection
}

// NewAuditRepo - creates a new Audit repository.
func NewAuditRepo(m *mongodb.MongoDB) *AuditRepo {
	return &AuditRepo{
		MongoDB:    m,
		collection: m.Database.Collection("audit"),
	}
}

// Create - inserts a new audit log record.
func (r *AuditRepo) Create(ctx context.Context, log entity.AuditLog) (*entity.AuditLog, error) {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	res, err := r.collection.InsertOne(ctx, log)
	if err != nil {
		return nil, fmt.Errorf("AuditRepo - Create - InsertOne: %w", err)
	}

	// set inserted ID back to entity
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		log.ID = oid.Hex()
	}

	return &log, nil
}

// GetAll - retrieves all audit logs.
func (r *AuditRepo) GetAll(ctx context.Context) ([]entity.AuditLog, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("AuditRepo - GetAll - Find: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []entity.AuditLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("AuditRepo - GetAll - Decode: %w", err)
	}

	if len(logs) == 0 {
		return nil, errors.New("AuditRepo - GetAll: no logs found")
	}

	return logs, nil
}

