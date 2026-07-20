package repositories

import "context"

type QueueRepository interface {
	AddObject(ctx context.Context, queueName string, data interface{}) error
}
