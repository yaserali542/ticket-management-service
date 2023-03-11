package models

import uuid "github.com/satori/go.uuid"

type BasicFields struct {
	ID           uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserName     string    `json:"user_name"`
	EmailAddress string    `json:"email_address"`
}

type BasicFieldsRequest struct {
	UserName string `json:"user_name"`
}
