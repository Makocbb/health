package providers

import (
	"context"
	"math/rand"
	"os"
	"time"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"
)

type emailRepository struct {
}

func NewEmailRepository() repositories.EmailRepository {
	return &emailRepository{}
}

// SendEmail is a mock SendGrid client.
// HEALTH_MOCK_RELIABLE:
//   - "1": always succeed, no sleep
//   - "0": always fail, no sleep
//   - unset/other: ~95% success with ~1s latency
func (r *emailRepository) SendEmail(ctx context.Context, req models.EmailRequest) error {
	switch os.Getenv("HEALTH_MOCK_RELIABLE") {
	case "1":
		return nil
	case "0":
		return models.ErrSendFailed
	}

	time.Sleep(time.Second)

	if rand.Float64() < 0.95 {
		return nil
	}

	return models.ErrSendFailed
}
