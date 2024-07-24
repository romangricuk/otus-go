package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type EventRepo struct {
	storage *MemoryStorage
	mu      sync.RWMutex
}

func (r *EventRepo) CreateEvent(_ context.Context, event storage.Event) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range r.storage.events {
		if e.StartTime.Before(event.EndTime) && event.StartTime.Before(e.EndTime) {
			return uuid.Nil, storage.ErrDateBusy
		}
	}
	event.ID = uuid.New()
	r.storage.events[event.ID] = event
	return event.ID, nil
}

func (r *EventRepo) UpdateEvent(_ context.Context, id uuid.UUID, event storage.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.storage.events[id]; !exists {
		return storage.ErrEventNotFound
	}
	for _, e := range r.storage.events {
		if e.ID != id && e.StartTime.Before(event.EndTime) && event.StartTime.Before(e.EndTime) {
			return storage.ErrDateBusy
		}
	}
	r.storage.events[id] = event
	return nil
}

func (r *EventRepo) DeleteEvent(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.storage.events[id]; !exists {
		return storage.ErrEventNotFound
	}
	delete(r.storage.events, id)
	return nil
}

func (r *EventRepo) GetEvent(_ context.Context, id uuid.UUID) (storage.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	event, exists := r.storage.events[id]
	if !exists {
		return storage.Event{}, storage.ErrEventNotFound
	}
	return event, nil
}

func (r *EventRepo) ListEvents(_ context.Context, start, end time.Time) ([]storage.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var events []storage.Event
	for _, event := range r.storage.events {
		if event.StartTime.After(start) && event.EndTime.Before(end) {
			events = append(events, event)
		}
	}
	return events, nil
}
