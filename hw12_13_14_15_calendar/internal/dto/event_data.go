package dto

import (
	"github.com/google/uuid"
	"time"
)

type EventData struct {
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	UserID      uuid.UUID
}
