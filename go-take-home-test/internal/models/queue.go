package models

import (
	"fmt"
	"time"
)

var ErrQueueSendFailed = fmt.Errorf("failed to send queue message")

// QueueOptions configures the mock local queue provider.
// Can be upgraded later with credentials for a real queue backend.
type QueueOptions struct {
	BaseURL        string
	MaxAttempts    int
	InitialBackoff time.Duration
}

func NewQueueOptions(inputOptions ...QueueOption) *QueueOptions {
	opts := &QueueOptions{
		BaseURL:        "http://localhost:3000",
		MaxAttempts:    3,
		InitialBackoff: time.Second,
	}
	for _, opt := range inputOptions {
		opt(opts)
	}
	return opts
}

// QueueOption is a functional option for configuring QueueOptions.
type QueueOption func(*QueueOptions)

// QueueWithBaseURL sets the base URL of this service (used to reach /workers).
func QueueWithBaseURL(baseURL string) QueueOption {
	return func(opts *QueueOptions) {
		opts.BaseURL = baseURL
	}
}

// QueueWithMaxAttempts sets how many times a failed delivery is retried (including the first try).
func QueueWithMaxAttempts(n int) QueueOption {
	return func(opts *QueueOptions) {
		opts.MaxAttempts = n
	}
}

// QueueWithInitialBackoff sets the delay before the first retry; doubles after each failure.
func QueueWithInitialBackoff(d time.Duration) QueueOption {
	return func(opts *QueueOptions) {
		opts.InitialBackoff = d
	}
}
