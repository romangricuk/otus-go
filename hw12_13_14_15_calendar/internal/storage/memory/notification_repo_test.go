package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/dto"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestNotificationRepo_CreateNotification(t *testing.T) {
	memStore := New()
	repo := memStore.NotificationRepository()

	notification := storage.Notification{
		EventID: uuid.New(),
		Time:    time.Now().Add(1 * time.Hour),
		Message: "Test Notification",
		Sent:    dto.NotificationOnWait,
	}

	id, err := repo.CreateNotification(context.Background(), notification)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestNotificationRepo_UpdateNotification(t *testing.T) {
	memStore := New()
	repo := memStore.NotificationRepository()

	notification := storage.Notification{
		EventID: uuid.New(),
		Time:    time.Now().Add(1 * time.Hour),
		Message: "Test Notification",
		Sent:    dto.NotificationOnWait,
	}

	id, _ := repo.CreateNotification(context.Background(), notification)
	notification.ID = id
	notification.Message = "Updated Message"

	err := repo.UpdateNotification(context.Background(), id, notification)
	assert.NoError(t, err)

	updatedNotification, err := repo.GetNotification(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Message", updatedNotification.Message)
}

func TestNotificationRepo_DeleteNotification(t *testing.T) {
	memStore := New()
	repo := memStore.NotificationRepository()

	notification := storage.Notification{
		EventID: uuid.New(),
		Time:    time.Now().Add(1 * time.Hour),
		Message: "Test Notification",
		Sent:    dto.NotificationOnWait,
	}

	id, _ := repo.CreateNotification(context.Background(), notification)

	err := repo.DeleteNotification(context.Background(), id)
	assert.NoError(t, err)

	_, err = repo.GetNotification(context.Background(), id)
	assert.ErrorIs(t, err, storage.ErrNotificationNotFound)
}

func TestNotificationRepo_ListNotifications(t *testing.T) {
	memStore := New()
	repo := memStore.NotificationRepository()

	notification1 := storage.Notification{
		EventID: uuid.New(),
		Time:    time.Now().Add(1 * time.Hour),
		Message: "Notification 1",
		Sent:    dto.NotificationOnWait,
	}

	notification2 := storage.Notification{
		EventID: uuid.New(),
		Time:    time.Now().Add(2 * time.Hour),
		Message: "Notification 2",
		Sent:    dto.NotificationOnWait,
	}

	_, _ = repo.CreateNotification(context.Background(), notification1)
	_, _ = repo.CreateNotification(context.Background(), notification2)

	notifications, err := repo.ListNotifications(context.Background(), time.Now(), time.Now().Add(3*time.Hour))
	assert.NoError(t, err)
	assert.Len(t, notifications, 2)
}
