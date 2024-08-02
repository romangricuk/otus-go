package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	UserID      uuid.UUID
}

type EventRepository interface {
	CreateEvent(ctx context.Context, event Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, event Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvent(ctx context.Context, id uuid.UUID) (Event, error)
	ListEvents(ctx context.Context, start, end time.Time) ([]Event, error)
}
