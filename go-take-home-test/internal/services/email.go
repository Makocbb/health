package services

import (
	"context"
	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"
)

type emailService struct {
	emailRepository repositories.EmailRepository
}

type EmailService interface {
	SendEmail(ctx context.Context, req models.EmailRequest) error
}

func NewEmailService(emailRepository repositories.EmailRepository) EmailService {
	return &emailService{emailRepository: emailRepository}
}

func (s *emailService) SendEmail(ctx context.Context, req models.EmailRequest) error {
	return s.emailRepository.SendEmail(ctx, req)
}
