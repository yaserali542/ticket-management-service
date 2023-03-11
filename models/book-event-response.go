package models

import uuid "github.com/satori/go.uuid"

type BookEventResponse struct {
	BookingId    uuid.UUID `json:"booking_id"`
	Signature    string    `json:"digital_signature"`
	EventDetails Events    `json:"event_details"`
	HashedData   string    `json:"hashed-data"`
	Otp          string    `json:"otp"`
}
