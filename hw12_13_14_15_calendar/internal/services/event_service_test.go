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

func TestEventService(t *testing.T) {
	ctx := context.Background()
	store := memorystorage.New()
	service := NewEventService(store.EventRepository())

	event := dto.EventData{
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	t.Run("CreateEvent", func(t *testing.T) {
		id, err := service.CreateEvent(ctx, event)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, id)
	})

	t.Run("GetEvent", func(t *testing.T) {
		id, err := service.CreateEvent(ctx, event)
		require.NoError(t, err)

		retrievedEvent, err := service.GetEvent(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, event.Title, retrievedEvent.Title)
		assert.Equal(t, event.Description, retrievedEvent.Description)
		assert.WithinDuration(t, event.StartTime, retrievedEvent.StartTime, time.Second)
		assert.WithinDuration(t, event.EndTime, retrievedEvent.EndTime, time.Second)
		assert.Equal(t, event.UserID, retrievedEvent.UserID)
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		id, err := service.CreateEvent(ctx, event)
		require.NoError(t, err)

		updatedEvent := dto.EventData{
			Title:       "Updated Event",
			Description: "This is an updated event",
			StartTime:   event.StartTime.Add(1 * time.Hour),
			EndTime:     event.EndTime.Add(1 * time.Hour),
			UserID:      event.UserID,
		}

		err = service.UpdateEvent(ctx, id, updatedEvent)
		require.NoError(t, err)

		retrievedEvent, err := service.GetEvent(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, updatedEvent.Title, retrievedEvent.Title)
		assert.Equal(t, updatedEvent.Description, retrievedEvent.Description)
		assert.WithinDuration(t, updatedEvent.StartTime, retrievedEvent.StartTime, time.Second)
		assert.WithinDuration(t, updatedEvent.EndTime, retrievedEvent.EndTime, time.Second)
		assert.Equal(t, updatedEvent.UserID, retrievedEvent.UserID)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		id, err := service.CreateEvent(ctx, event)
		require.NoError(t, err)

		err = service.DeleteEvent(ctx, id)
		require.NoError(t, err)

		_, err = service.GetEvent(ctx, id)
		assert.Error(t, err)
	})

	t.Run("ListEvents", func(t *testing.T) {
		start := time.Now()
		end := start.Add(24 * time.Hour)

		events, err := service.ListEvents(ctx, start, end)
		require.NoError(t, err)
		assert.NotEmpty(t, events)
	})
}
