package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID      uuid.UUID `json:"id"`
	EventID uuid.UUID `json:"eventId"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Sent    bool      `json:"sent"`
}

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification Notification) (uuid.UUID, error)
	UpdateNotification(ctx context.Context, id uuid.UUID, notification Notification) error
	DeleteNotification(ctx context.Context, id uuid.UUID) error
	GetNotification(ctx context.Context, id uuid.UUID) (Notification, error)
	ListNotifications(ctx context.Context, start, end time.Time) ([]Notification, error)
}
