package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/dto"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationService(t *testing.T) {
	ctx := context.Background()
	store := memorystorage.New()
	service := NewNotificationService(store.NotificationRepository())

	notification := dto.NotificationData{
		EventID: uuid.New(),
		UserID:  uuid.New(),
		Time:    time.Now(),
		Message: "Test Notification",
		Sent:    false,
	}

	t.Run("CreateNotification", func(t *testing.T) {
		id, err := service.CreateNotification(ctx, notification)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, id)
	})

	t.Run("GetNotification", func(t *testing.T) {
		id, err := service.CreateNotification(ctx, notification)
		require.NoError(t, err)

		retrievedNotification, err := service.GetNotification(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, notification.EventID, retrievedNotification.EventID)
		assert.Equal(t, notification.UserID, retrievedNotification.UserID)
		assert.WithinDuration(t, notification.Time, retrievedNotification.Time, time.Second)
		assert.Equal(t, notification.Message, retrievedNotification.Message)
		assert.Equal(t, notification.Sent, retrievedNotification.Sent)
	})

	t.Run("UpdateNotification", func(t *testing.T) {
		id, err := service.CreateNotification(ctx, notification)
		require.NoError(t, err)

		updatedNotification := dto.NotificationData{
			EventID: notification.EventID,
			UserID:  notification.UserID,
			Time:    notification.Time.Add(1 * time.Hour),
			Message: "Updated Notification",
			Sent:    true,
		}

		err = service.UpdateNotification(ctx, id, updatedNotification)
		require.NoError(t, err)

		retrievedNotification, err := service.GetNotification(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, updatedNotification.EventID, retrievedNotification.EventID)
		assert.Equal(t, updatedNotification.UserID, retrievedNotification.UserID)
		assert.WithinDuration(t, updatedNotification.Time, retrievedNotification.Time, time.Second)
		assert.Equal(t, updatedNotification.Message, retrievedNotification.Message)
		assert.Equal(t, updatedNotification.Sent, retrievedNotification.Sent)
	})

	t.Run("DeleteNotification", func(t *testing.T) {
		id, err := service.CreateNotification(ctx, notification)
		require.NoError(t, err)

		err = service.DeleteNotification(ctx, id)
		require.NoError(t, err)

		_, err = service.GetNotification(ctx, id)
		assert.Error(t, err)
	})

	t.Run("ListNotifications", func(t *testing.T) {
		start := time.Now()
		end := start.Add(24 * time.Hour)

		notifications, err := service.ListNotifications(ctx, start, end)
		require.NoError(t, err)
		assert.NotEmpty(t, notifications)
	})
}
