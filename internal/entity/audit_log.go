// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import (
	"time"
)

// AuditLog - represents an audit log entry in MongoDB.
type AuditLog struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	Event         string    `bson:"event" json:"event"`
	Payload       any       `bson:"payload,omitempty" json:"payload,omitempty"`
	Timestamp     time.Time `bson:"timestamp" json:"timestamp"`
}
