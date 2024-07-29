package dto

import (
	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/api"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type NotificationData struct {
	ID      uuid.UUID `json:"id,omitempty"`
	EventID uuid.UUID `json:"eventId"`
	UserID  uuid.UUID `json:"userId"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Sent    bool      `json:"sent"`
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

func ToApiNotification(notification NotificationData) *api.Notification {
	return &api.Notification{
		Id:      notification.ID.String(),
		EventId: notification.EventID.String(),
		UserId:  notification.UserID.String(),
		Time:    timestamppb.New(notification.Time),
		Message: notification.Message,
		Sent:    notification.Sent,
	}
}

func FromApiNotification(notification *api.Notification) NotificationData {
	return NotificationData{
		ID:      uuid.MustParse(notification.GetId()),
		EventID: uuid.MustParse(notification.GetEventId()),
		UserID:  uuid.MustParse(notification.GetUserId()),
		Time:    notification.GetTime().AsTime(),
		Message: notification.GetMessage(),
		Sent:    notification.GetSent(),
	}
}
