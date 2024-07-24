package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	UserID      uuid.UUID `json:"userId"`
}

type EventRepository interface {
	CreateEvent(ctx context.Context, event Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, event Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvent(ctx context.Context, id uuid.UUID) (Event, error)
	ListEvents(ctx context.Context, start, end time.Time) ([]Event, error)
}
