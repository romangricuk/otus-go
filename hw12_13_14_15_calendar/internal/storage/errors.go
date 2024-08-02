package storage

import "errors"

var (
	ErrEventNotFound        = errors.New("event not found")
	ErrNotificationNotFound = errors.New("notification not found")
)
