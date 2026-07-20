package services

import (
	"context"

	"go-take-home-test/internal/repositories"
)

type queueService struct {
	repo repositories.QueueRepository
}

type QueueService interface {
	AddObject(ctx context.Context, queueName string, data any) error
}

func NewQueueService(repo repositories.QueueRepository) QueueService {
	return &queueService{repo: repo}
}

func (s *queueService) AddObject(ctx context.Context, queueName string, data any) error {
	return s.repo.AddObject(ctx, queueName, data)
}
