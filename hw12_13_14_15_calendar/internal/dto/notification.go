package dto

import (
	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type NotificationData struct {
	ID      uuid.UUID `json:"id"`
	EventID uuid.UUID `json:"eventId"`
	UserID  uuid.UUID `json:"userId"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Sent    bool      `json:"sent"`
}

func ToStorageNotification(notification NotificationData) storage.Notification {
	return storage.Notification{
		ID:      notification.ID,
		EventID: notification.EventID,
		UserID:  notification.UserID,
		Time:    notification.Time,
		Message: notification.Message,
		Sent:    notification.Sent,
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
