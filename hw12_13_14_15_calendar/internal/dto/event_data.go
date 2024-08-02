package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/api"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventData struct {
	ID          uuid.UUID `json:"id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title       string    `json:"title" example:"Event title"`
	Description string    `json:"description" example:"Event description"`
	StartTime   time.Time `json:"startTime" example:"2024-07-02T00:00:00Z"`
	EndTime     time.Time `json:"endTime" example:"2024-07-02T00:00:00Z"`
	UserID      uuid.UUID `json:"userId" example:"123e4567-e89b-12d3-a456-426614174000"`
}

func ToStorageEvent(data EventData) storage.Event {
	return storage.Event{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		StartTime:   data.StartTime,
		EndTime:     data.EndTime,
		UserID:      data.UserID,
	}
}

func FromStorageEvent(event storage.Event) EventData {
	return EventData{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		UserID:      event.UserID,
	}
}

func ToAPIEvent(event EventData) *api.Event {
	return &api.Event{
		Id:          event.ID.String(),
		Title:       event.Title,
		Description: event.Description,
		StartTime:   timestamppb.New(event.StartTime),
		EndTime:     timestamppb.New(event.EndTime),
		UserId:      event.UserID.String(),
	}
}

func FromAPIEvent(event *api.Event) EventData {
	return EventData{
		ID:          uuid.MustParse(event.GetId()),
		Title:       event.GetTitle(),
		Description: event.GetDescription(),
		StartTime:   event.GetStartTime().AsTime(),
		EndTime:     event.GetEndTime().AsTime(),
		UserID:      uuid.MustParse(event.GetUserId()),
	}
}
