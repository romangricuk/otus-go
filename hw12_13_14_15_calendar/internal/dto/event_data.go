package dto

import (
	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type EventData struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	UserID      uuid.UUID `json:"userId"`
}

func ToStorageEvent(data EventData) storage.Event {
	return storage.Event{
		Title:       data.Title,
		Description: data.Description,
		StartTime:   data.StartTime,
		EndTime:     data.EndTime,
		UserID:      data.UserID,
	}
}

func FromStorageEvent(event storage.Event) EventData {
	return EventData{
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		UserID:      event.UserID,
	}
}
