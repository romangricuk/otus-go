package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/api"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	NotificationOnWait  = "wait"
	NotificationOnQueue = "on-queue"
	NotificationSent    = "sent"
)

type NotificationData struct {
	ID      uuid.UUID `json:"id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	EventID uuid.UUID `json:"eventId" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID  uuid.UUID `json:"userId" example:"123e4567-e89b-12d3-a456-426614174000"`
	Time    time.Time `json:"time" example:"2024-07-02T00:00:00Z"`
	Message string    `json:"message" example:"Notification message"`
	Sent    string    `json:"sent" example:"wait, on-queue, sent"`
}

func ToStorageNotification(data NotificationData) storage.Notification {
	return storage.Notification{
		ID:      data.ID,
		EventID: data.EventID,
		UserID:  data.UserID,
		Time:    data.Time,
		Message: data.Message,
		Sent:    data.Sent,
	}
}

func FromStorageNotification(notification storage.Notification) NotificationData {
	return NotificationData{
		ID:      notification.ID,
		EventID: notification.EventID,
		UserID:  notification.UserID,
		Time:    notification.Time,
		Message: notification.Message,
		Sent:    notification.Sent,
	}
}

func ToAPINotification(notification NotificationData) *api.Notification {
	return &api.Notification{
		Id:      notification.ID.String(),
		EventId: notification.EventID.String(),
		UserId:  notification.UserID.String(),
		Time:    timestamppb.New(notification.Time),
		Message: notification.Message,
		Sent:    notification.Sent,
	}
}

func FromAPINotification(notification *api.Notification) NotificationData {
	return NotificationData{
		ID:      uuid.MustParse(notification.GetId()),
		EventID: uuid.MustParse(notification.GetEventId()),
		UserID:  uuid.MustParse(notification.GetUserId()),
		Time:    notification.GetTime().AsTime(),
		Message: notification.GetMessage(),
		Sent:    notification.GetSent(),
	}
}
