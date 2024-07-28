package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type NotificationService interface {
	CreateNotification(
		ctx context.Context,
		eventID uuid.UUID,
		userID uuid.UUID,
		time time.Time,
		message string,
		sent bool,
	) (uuid.UUID, error)
	UpdateNotification(
		ctx context.Context,
		id uuid.UUID,
		eventID uuid.UUID,
		userID uuid.UUID,
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

func NewNotificationService(repo storage.NotificationRepository) *NotificationServiceImpl {
	return &NotificationServiceImpl{storage: repo}
}

func (s *NotificationServiceImpl) CreateNotification(
	ctx context.Context,
	eventID uuid.UUID,
	userID uuid.UUID,
	time time.Time,
	message string,
	sent bool,
) (uuid.UUID, error) {
	notification := storage.Notification{
		ID:      uuid.New(),
		EventID: eventID,
		UserID:  userID,
		Time:    time,
		Message: message,
		Sent:    sent,
	}
	return s.storage.CreateNotification(ctx, notification)
}

func (s *NotificationServiceImpl) UpdateNotification(
	ctx context.Context,
	id uuid.UUID,
	eventID uuid.UUID,
	userID uuid.UUID,
	time time.Time,
	message string,
	sent bool,
) error {
	notification := storage.Notification{
		ID:      id,
		EventID: eventID,
		UserID:  userID,
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
	start time.Time,
	end time.Time,
) ([]storage.Notification, error) {
	return s.storage.ListNotifications(ctx, start, end)
}
