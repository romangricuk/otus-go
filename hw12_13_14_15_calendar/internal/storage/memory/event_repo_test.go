package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestEventRepo_CreateEvent(t *testing.T) {
	memStore := New()
	repo := memStore.EventRepository()

	event := storage.Event{
		Title:       "Test Event",
		Description: "Description",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	id, err := repo.CreateEvent(context.Background(), event)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestEventRepo_UpdateEvent(t *testing.T) {
	memStore := New()
	repo := memStore.EventRepository()

	event := storage.Event{
		Title:       "Test Event",
		Description: "Description",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	id, _ := repo.CreateEvent(context.Background(), event)
	event.ID = id
	event.Title = "Updated Title"

	err := repo.UpdateEvent(context.Background(), id, event)
	assert.NoError(t, err)

	updatedEvent, err := repo.GetEvent(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedEvent.Title)
}

func TestEventRepo_DeleteEvent(t *testing.T) {
	memStore := New()
	repo := memStore.EventRepository()

	event := storage.Event{
		Title:       "Test Event",
		Description: "Description",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	id, _ := repo.CreateEvent(context.Background(), event)

	err := repo.DeleteEvent(context.Background(), id)
	assert.NoError(t, err)

	_, err = repo.GetEvent(context.Background(), id)
	assert.Error(t, err)
}

func TestEventRepo_ListEvents(t *testing.T) {
	memStore := New()
	repo := memStore.EventRepository()

	event1 := storage.Event{
		Title:       "Event 1",
		Description: "Description 1",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	event2 := storage.Event{
		Title:       "Event 2",
		Description: "Description 2",
		StartTime:   time.Now().Add(2 * time.Hour),
		EndTime:     time.Now().Add(3 * time.Hour),
		UserID:      uuid.New(),
	}

	_, _ = repo.CreateEvent(context.Background(), event1)
	_, _ = repo.CreateEvent(context.Background(), event2)

	events, err := repo.ListEvents(context.Background(), time.Now().Add(-time.Minute), time.Now().Add(4*time.Hour))
	assert.NoError(t, err)
	assert.Len(t, events, 2)
}
