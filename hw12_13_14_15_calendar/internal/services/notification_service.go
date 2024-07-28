package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, eventID uuid.UUID, time time.Time, message string) (uuid.UUID, error)
	UpdateNotification(
		ctx context.Context,
		id uuid.UUID,
		eventID uuid.UUID,
		time time.Time,
		message string,
		sent bool,
	) error
	DeleteNotification(ctx context.Context, id uuid.UUID) error
	GetNotification(ctx context.Context, id uuid.UUID) (storage.Notification, error)
	ListNotifications(ctx context.Context, start, end time.Time) ([]storage.Notification, error)
}

type NotificationServiceImpl struct {
	storage storage.NotificationRepository
}

func NewNotificationService(storage storage.NotificationRepository) *NotificationServiceImpl {
	return &NotificationServiceImpl{storage: storage}
}

func (s *NotificationServiceImpl) CreateNotification(
	ctx context.Context,
	eventID uuid.UUID,
	time time.Time,
	message string,
) (uuid.UUID, error) {
	notification := storage.Notification{
		ID:      uuid.New(),
		EventID: eventID,
		Time:    time,
		Message: message,
		Sent:    false,
	}
	return s.storage.CreateNotification(ctx, notification)
}

func (s *NotificationServiceImpl) UpdateNotification(
	ctx context.Context,
	id uuid.UUID,
	eventID uuid.UUID,
	time time.Time,
	message string,
	sent bool,
) error {
	notification := storage.Notification{
		ID:      id,
		EventID: eventID,
		Time:    time,
		Message: message,
		Sent:    sent,
	}
	return s.storage.UpdateNotification(ctx, id, notification)
}

func (s *NotificationServiceImpl) DeleteNotification(ctx context.Context, id uuid.UUID) error {
	return s.storage.DeleteNotification(ctx, id)
}

func (s *NotificationServiceImpl) GetNotification(ctx context.Context, id uuid.UUID) (storage.Notification, error) {
	return s.storage.GetNotification(ctx, id)
}

func (s *NotificationServiceImpl) ListNotifications(
	ctx context.Context,
	start,
	end time.Time,
) ([]storage.Notification, error) {
	return s.storage.ListNotifications(ctx, start, end)
}
