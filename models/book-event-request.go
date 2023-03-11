package models

import uuid "github.com/satori/go.uuid"

type BookEventRequest struct {
	EventId  uuid.UUID `json:"event_id"`
	UserName string    `json:"user_name"`
}
