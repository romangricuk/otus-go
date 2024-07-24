package storage

import "errors"

var (
	ErrDateBusy              = errors.New("the specified date and time are already occupied by another event")
	ErrEventNotFound         = errors.New("event not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrNotificationNotFound  = errors.New("notification not found")
	ErrNotificationTimePast  = errors.New("notification time has already passed")
	ErrNotificationDuplicate = errors.New("duplicate notification for the same event and time")
)
