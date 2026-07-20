package repositories

import (
	"context"
	"go-take-home-test/internal/models"
)

type EmailRepository interface {
	SendEmail(ctx context.Context, req models.EmailRequest) error
}
