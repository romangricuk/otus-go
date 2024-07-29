package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type NotificationRepo struct {
	notifications map[uuid.UUID]storage.Notification
	mu            sync.RWMutex
}

func (r *NotificationRepo) CreateNotification(_ context.Context, notification storage.Notification) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	notification.ID = uuid.New()
	r.notifications[notification.ID] = notification
	return notification.ID, nil
}

func (r *NotificationRepo) UpdateNotification(_ context.Context, id uuid.UUID, notification storage.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.notifications[id]; !exists {
		return storage.ErrNotificationNotFound
	}
	notification.ID = id
	r.notifications[id] = notification
	return nil
}

func (r *NotificationRepo) DeleteNotification(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.notifications[id]; !exists {
		return storage.ErrNotificationNotFound
	}
	delete(r.notifications, id)
	return nil
}

func (r *NotificationRepo) GetNotification(_ context.Context, id uuid.UUID) (storage.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	notification, exists := r.notifications[id]
	if !exists {
		return storage.Notification{}, storage.ErrNotificationNotFound
	}
	return notification, nil
}

func (r *NotificationRepo) ListNotifications(_ context.Context, start, end time.Time) ([]storage.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var notifications []storage.Notification
	for _, notification := range r.notifications {
		if notification.Time.After(start) && notification.Time.Before(end) {
			notifications = append(notifications, notification)
		}
	}
	return notifications, nil
}
