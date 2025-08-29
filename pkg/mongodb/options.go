package mongodb

import "time"

// Option -.
type Option func(*MongoDB)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *MongoDB) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *MongoDB) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *MongoDB) {
		c.connTimeout = timeout
	}
}
