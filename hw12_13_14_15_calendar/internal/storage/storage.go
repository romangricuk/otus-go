package storage

import "context"

type Storage interface {
	Connect(ctx context.Context) error
	Close() error
	HealthCheck(ctx context.Context) error
	EventRepository() EventRepository
	NotificationRepository() NotificationRepository
}
