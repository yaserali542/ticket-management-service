package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Booking struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"booking_id"`
	UserId    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	EventId   uuid.UUID `gorm:"type:uuid;not null" json:"event_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
