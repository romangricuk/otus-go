package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventRepository is a mock implementation of the EventRepository interface.
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) CreateEvent(ctx context.Context, event storage.Event) (uuid.UUID, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockEventRepository) UpdateEvent(ctx context.Context, id uuid.UUID, event storage.Event) error {
	args := m.Called(ctx, id, event)
	return args.Error(0)
}

func (m *MockEventRepository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventRepository) GetEvent(ctx context.Context, id uuid.UUID) (storage.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(storage.Event), args.Error(1)
}

func (m *MockEventRepository) ListEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockEventRepository) ListUpcomingEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func TestEventService(t *testing.T) {
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)

	ctx := context.Background()
	eventID := uuid.New()
	userID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(2 * time.Hour)
	event := storage.Event{
		ID:          eventID,
		Title:       "Test Event",
		Description: "Test Description",
		StartTime:   startTime,
		EndTime:     endTime,
		UserID:      userID,
	}

	t.Run("CreateEvent", func(t *testing.T) {
		mockRepo.On("CreateEvent", ctx, mock.AnythingOfType("storage.Event")).Return(eventID, nil).Once()

		id, err := service.CreateEvent(ctx, event.Title, event.Description, startTime, endTime, userID)
		assert.NoError(t, err)
		assert.Equal(t, eventID, id)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		mockRepo.On("UpdateEvent", ctx, eventID, mock.AnythingOfType("storage.Event")).Return(nil).Once()

		err := service.UpdateEvent(ctx, eventID, event.Title, event.Description, startTime, endTime, userID)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		mockRepo.On("DeleteEvent", ctx, eventID).Return(nil).Once()

		err := service.DeleteEvent(ctx, eventID)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetEvent", func(t *testing.T) {
		mockRepo.On("GetEvent", ctx, eventID).Return(event, nil).Once()

		result, err := service.GetEvent(ctx, eventID)
		assert.NoError(t, err)
		assert.Equal(t, event, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ListEvents", func(t *testing.T) {
		events := []storage.Event{event}
		mockRepo.On("ListEvents", ctx, startTime, endTime).Return(events, nil).Once()

		result, err := service.ListEvents(ctx, startTime, endTime)
		assert.NoError(t, err)
		assert.Equal(t, events, result)

		mockRepo.AssertExpectations(t)
	})
}
