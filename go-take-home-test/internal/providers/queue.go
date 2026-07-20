package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"
)

type queueEntry struct {
	QueuedAt time.Time
	Data     any
}

type queueRepository struct {
	opts   *models.QueueOptions
	client *http.Client

	mu    sync.Mutex
	cache map[string][]queueEntry
}

func NewQueueRepository(opts *models.QueueOptions) repositories.QueueRepository {
	if opts == nil {
		opts = models.NewQueueOptions()
	}
	return &queueRepository{
		opts:   opts,
		client: &http.Client{Timeout: 10 * time.Second},
		cache:  make(map[string][]queueEntry),
	}
}

// AddObject caches the payload locally, then POSTs it to /workers/<queueName>
// with exponential backoff retries on failure.
func (r *queueRepository) AddObject(ctx context.Context, queueName string, data any) error {
	r.mu.Lock()
	r.cache[queueName] = append(r.cache[queueName], queueEntry{
		QueuedAt: time.Now(),
		Data:     data,
	})
	r.mu.Unlock()

	if err := r.sendWithRetry(ctx, queueName, data); err != nil {
		slog.Error("failed to send queue object to worker", "error", err, "queue", queueName)
		return fmt.Errorf("%w: %v", models.ErrQueueSendFailed, err)
	}
	return nil
}

func (r *queueRepository) sendWithRetry(ctx context.Context, queueName string, data any) error {
	var lastErr error
	backoff := r.opts.InitialBackoff

	for attempt := 1; attempt <= r.opts.MaxAttempts; attempt++ {
		lastErr = r.send(ctx, queueName, data)
		if lastErr == nil {
			return nil
		}

		if attempt == r.opts.MaxAttempts {
			break
		}

		slog.Warn("queue delivery failed, retrying",
			"error", lastErr,
			"queue", queueName,
			"attempt", attempt,
			"backoff", backoff,
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
		backoff *= 2
	}

	return lastErr
}

func (r *queueRepository) send(ctx context.Context, queueName string, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	endpoint, err := url.JoinPath(r.opts.BaseURL, "workers", queueName)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("worker returned status %d", resp.StatusCode)
	}
	return nil
}
