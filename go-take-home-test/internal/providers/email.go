package providers

import (
	"context"
	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"
	"math/rand"
	"time"
)

type emailRepository struct {
}

func NewEmailRepository() repositories.EmailRepository {
	return &emailRepository{}
}

// SendEmail is a mock SendGrid client.
// It succeeds ~95% of the time and always sleeps ~1s to simulate latency.
func (r *emailRepository) SendEmail(ctx context.Context, req models.EmailRequest) error {
	time.Sleep(time.Second)

	if rand.Float64() < 0.95 {
		return nil
	}

	return models.ErrSendFailed
}
