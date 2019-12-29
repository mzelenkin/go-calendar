package entities

import (
	"time"
)

// Event событие в календаре
type Event struct {
	ID          EventID
	Title       string
	Start       time.Time
	End         time.Time
	Description string
}
