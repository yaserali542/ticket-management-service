package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Events struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	EventName  string    `json:"event_name"`
	EventImage []byte    `json:"event_image"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Venue      string    `json:"venue"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
