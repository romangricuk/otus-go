package memorystorage

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type MemoryStorage struct {
	eventRepo        *EventRepo
	notificationRepo *NotificationRepo
}

func New() *MemoryStorage {
	store := &MemoryStorage{
		eventRepo:        &EventRepo{events: make(map[uuid.UUID]storage.Event), mu: sync.RWMutex{}},
		notificationRepo: &NotificationRepo{notifications: make(map[uuid.UUID]storage.Notification), mu: sync.RWMutex{}},
	}
	return store
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
	return s.eventRepo
}

func (s *MemoryStorage) NotificationRepository() storage.NotificationRepository {
	return s.notificationRepo
}

func (s *MemoryStorage) HealthCheck(context.Context) error {
	return nil
}
