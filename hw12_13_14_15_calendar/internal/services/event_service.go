package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type EventService interface {
	CreateEvent(
		ctx context.Context,
		title string,
		description string,
		startTime time.Time,
		endTime time.Time,
		userID uuid.UUID,
	) (uuid.UUID, error)

	UpdateEvent(
		ctx context.Context,
		id uuid.UUID,
		title string,
		description string,
		startTime time.Time,
		endTime time.Time,
		userID uuid.UUID,
	) error

	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvent(ctx context.Context, id uuid.UUID) (storage.Event, error)
	ListEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error)
}

type EventServiceImpl struct {
	repo storage.EventRepository
}

func NewEventService(repo storage.EventRepository) *EventServiceImpl {
	return &EventServiceImpl{repo: repo}
}

func (s *EventServiceImpl) CreateEvent(
	ctx context.Context,
	title string,
	description string,
	startTime time.Time,
	endTime time.Time,
	userID uuid.UUID,
) (uuid.UUID, error) {
	event := storage.Event{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		UserID:      userID,
	}
	return s.repo.CreateEvent(ctx, event)
}

func (s *EventServiceImpl) UpdateEvent(
	ctx context.Context,
	id uuid.UUID,
	title string,
	description string,
	startTime time.Time,
	endTime time.Time,
	userID uuid.UUID,
) error {
	event := storage.Event{
		ID:          id,
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		UserID:      userID,
	}
	return s.repo.UpdateEvent(ctx, id, event)
}

func (s *EventServiceImpl) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteEvent(ctx, id)
}

func (s *EventServiceImpl) GetEvent(ctx context.Context, id uuid.UUID) (storage.Event, error) {
	return s.repo.GetEvent(ctx, id)
}

func (s *EventServiceImpl) ListEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	return s.repo.ListEvents(ctx, start, end)
}
