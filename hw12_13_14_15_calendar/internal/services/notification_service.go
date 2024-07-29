package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/dto"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, notification dto.NotificationData) (uuid.UUID, error)
	UpdateNotification(ctx context.Context, id uuid.UUID, notification dto.NotificationData) error
	DeleteNotification(ctx context.Context, id uuid.UUID) error
	GetNotification(ctx context.Context, id uuid.UUID) (dto.NotificationData, error)
	ListNotifications(ctx context.Context, start, end time.Time) ([]dto.NotificationData, error)
}

type NotificationServiceImpl struct {
	repo storage.NotificationRepository
}

func NewNotificationService(repo storage.NotificationRepository) NotificationService {
	return &NotificationServiceImpl{repo: repo}
}

func (s *NotificationServiceImpl) CreateNotification(ctx context.Context, notification dto.NotificationData) (uuid.UUID, error) {
	storageNotification := dto.ToStorageNotification(notification)
	storageNotification.ID = uuid.New()
	return s.repo.CreateNotification(ctx, storageNotification)
}

func (s *NotificationServiceImpl) UpdateNotification(ctx context.Context, id uuid.UUID, notification dto.NotificationData) error {
	storageNotification := dto.ToStorageNotification(notification)
	return s.repo.UpdateNotification(ctx, id, storageNotification)
}

func (s *NotificationServiceImpl) DeleteNotification(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteNotification(ctx, id)
}

func (s *NotificationServiceImpl) GetNotification(ctx context.Context, id uuid.UUID) (dto.NotificationData, error) {
	storageNotification, err := s.repo.GetNotification(ctx, id)
	if err != nil {
		return dto.NotificationData{}, err
	}
	return dto.FromStorageNotification(storageNotification), nil
}

func (s *NotificationServiceImpl) ListNotifications(ctx context.Context, start, end time.Time) ([]dto.NotificationData, error) {
	storageNotifications, err := s.repo.ListNotifications(ctx, start, end)
	if err != nil {
		return nil, err
	}
	notifications := make([]dto.NotificationData, len(storageNotifications))
	for i, storageNotification := range storageNotifications {
		notifications[i] = dto.FromStorageNotification(storageNotification)
	}
	return notifications, nil
}
