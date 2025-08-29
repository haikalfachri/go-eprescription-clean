// Package mongodb implements mongodb connection.
package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// MongoDB -.
type MongoDB struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Client   *mongo.Client
	Database *mongo.Database
}

// New -.
func New(url string, opts ...Option) (*MongoDB, error) {
	m := &MongoDB{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(m)
	}

	clientOpts := options.Client().
		ApplyURI(url).
		SetMaxPoolSize(uint64(m.maxPoolSize))

	var err error
	for m.connAttempts > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), m.connTimeout)
		defer cancel()

		m.Client, err = mongo.Connect(ctx, clientOpts)
		if err == nil {
			// Ping to check connection
			err = m.Client.Ping(ctx, nil)
			if err == nil {
				m.Database = m.Client.Database(os.Getenv("MG_DB_NAME"))
				break
			}
		}

		log.Printf("MongoDB is trying to connect, attempts left: %d", m.connAttempts)
		time.Sleep(m.connTimeout)
		m.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("mongodb - NewMongoDB - connAttempts == 0: %w", err)
	}

	return m, nil
}

// Close -.
func (m *MongoDB) Close() {
	if m.Client != nil {
		_ = m.Client.Disconnect(context.Background())
	}
}
