package memorystorage

import (
	"context"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type MemoryStorage struct {
	events        map[uuid.UUID]storage.Event
	notifications map[uuid.UUID]storage.Notification
}

func New() *MemoryStorage {
	return &MemoryStorage{
		events:        make(map[uuid.UUID]storage.Event),
		notifications: make(map[uuid.UUID]storage.Notification),
	}
}

func (s *MemoryStorage) Connect(context.Context) error {
	// No connection required for in-memory storage
	return nil
}

func (s *MemoryStorage) Close() error {
	// No connection to close for in-memory storage
	return nil
}

func (s *MemoryStorage) EventRepository() storage.EventRepository {
	return &EventRepo{storage: s}
}

func (s *MemoryStorage) NotificationRepository() storage.NotificationRepository {
	return &NotificationRepo{storage: s}
}

func (s *MemoryStorage) HealthCheck(context.Context) error {
	return nil
}
