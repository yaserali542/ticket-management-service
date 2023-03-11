package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type BookingData struct {
	ID               uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"-"`
	BookingId        uuid.UUID `gorm:"type:uuid;not null;index" json:"booking_id"`
	QrEncryptedData  []byte    `gorm:"not null" json:"-"`
	DigitalSignature string    `gorm:"not null" json:"signature"`
	Otp              string    `gorm:"not null"`
	CreatedAt        time.Time `json:"booking_date"`
	UpdatedAt        time.Time
}
