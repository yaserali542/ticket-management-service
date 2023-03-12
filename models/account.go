package models

import uuid "github.com/satori/go.uuid"

type Account struct { // table name: parents
	ID             uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"-"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	ProfilePicture []byte    `json:"profile_picture"`
	UserName       string    `json:"user_name"`
	EmailAddress   string    `json:"email_address"`
	HashedPassword []byte    `json:"-"`
	Salt           []byte    `json:"-"`
	Role           string    `gorm:"default:'customer'" json:"role"`
}
