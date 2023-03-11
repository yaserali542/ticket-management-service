package repository

import (
	"errors"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/yaserali542/ticket-management-service/models"
)

type Repository struct {
	Db *gorm.DB
}

func (r *Repository) GetEventDetails(eventId uuid.UUID) (*models.Events, error) {

	var event models.Events
	db := r.Db.First(&event, "id =?", eventId)
	if db.Error != nil {
		return nil, db.Error
	}
	return &event, nil
}

func (r *Repository) GetBookingDetails(userId uuid.UUID) (*models.Booking, bool, error) {

	var booking models.Booking

	db := r.Db.First(&booking, "user_id =?", userId)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, true, nil
		} else {
			return nil, true, db.Error
		}
	}

	return &booking, false, nil
}

func (r *Repository) GetEvents() ([]models.Events, error) {

	var events []models.Events

	if err := r.Db.Limit(10).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *Repository) CreateBooking(booking *models.Booking) error {
	return r.Db.Create(booking).Error
}
func (r *Repository) CreateBookingData(bookingData *models.BookingData) error {
	return r.Db.Create(bookingData).Error
}
