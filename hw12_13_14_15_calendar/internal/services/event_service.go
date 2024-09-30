package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/dto"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventService interface {
	CreateEvent(ctx context.Context, event dto.EventData) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, event dto.EventData) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvent(ctx context.Context, id uuid.UUID) (dto.EventData, error)
	ListEvents(ctx context.Context, start, end time.Time) ([]dto.EventData, error)
}

type EventServiceImpl struct {
	repo storage.EventRepository
}

func NewEventService(store storage.Storage) EventService {
	return &EventServiceImpl{repo: store.EventRepository()}
}

func (s *EventServiceImpl) CreateEvent(ctx context.Context, event dto.EventData) (uuid.UUID, error) {
	storageEvent := dto.ToStorageEvent(event)
	storageEvent.ID = uuid.New()

	newStart := storageEvent.StartTime
	newEnd := storageEvent.EndTime

	if newStart.After(newEnd) || newStart.Equal(newEnd) {
		return uuid.Nil, status.Error(codes.InvalidArgument, "the beginning of events must be before the end")
	}

	return s.repo.CreateEvent(ctx, storageEvent)
}

func (s *EventServiceImpl) UpdateEvent(ctx context.Context, id uuid.UUID, event dto.EventData) error {
	storageEvent := dto.ToStorageEvent(event)
	return s.repo.UpdateEvent(ctx, id, storageEvent)
}

func (s *EventServiceImpl) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteEvent(ctx, id)
}

func (s *EventServiceImpl) GetEvent(ctx context.Context, id uuid.UUID) (dto.EventData, error) {
	storageEvent, err := s.repo.GetEvent(ctx, id)
	if err != nil {
		return dto.EventData{}, err
	}
	return dto.FromStorageEvent(storageEvent), nil
}

func (s *EventServiceImpl) ListEvents(ctx context.Context, start, end time.Time) ([]dto.EventData, error) {
	storageEvents, err := s.repo.ListEvents(ctx, start, end)
	if err != nil {
		return nil, err
	}
	events := make([]dto.EventData, len(storageEvents))
	for i, storageEvent := range storageEvents {
		events[i] = dto.FromStorageEvent(storageEvent)
	}
	return events, nil
}
