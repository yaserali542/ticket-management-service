package repository

import (
	"encoding/base64"
	"errors"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/yaserali542/ticket-management-service/encryption"
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

func (r *Repository) GetBookingDetails(userId uuid.UUID, userName string) ([]models.BookEventResponse, bool, error) {

	var bookings []models.Booking
	//response := make([]models.BookEventResponse, 1)
	db := r.Db.Find(&bookings, "user_id =?", userId)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, true, nil
		} else {
			return nil, true, db.Error
		}
	}
	response := make([]models.BookEventResponse, len(bookings))
	for _, v := range bookings {
		var bookingDetail models.BookingData
		if err := r.Db.Find(&bookingDetail, "booking_id =?", v.ID).Error; err != nil {
			return nil, true, err
		}
		var event models.Events
		if err := r.Db.Find(&event, " id=?", v.EventId).Error; err != nil {
			return nil, true, err
		}

		hashedData := encryption.CalculateHash(bookingDetail.QrEncryptedData)

		response = append(response, models.BookEventResponse{
			EventDetails: models.Events{
				ID:         event.ID,
				EventImage: event.EventImage,
				StartDate:  event.StartDate,
				EventName:  event.EventName,
				EndDate:    event.EndDate,
				Venue:      event.Venue,
			},
			BookingId:  v.ID,
			Signature:  bookingDetail.DigitalSignature,
			Otp:        bookingDetail.Otp,
			HashedData: base64.StdEncoding.EncodeToString(hashedData),
			UserId:     userId.String(),
			UserName:   userName,
		})

	}

	return response, false, nil
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

func (r *Repository) ViewBooking() ([]models.Events, error) {

	var events []models.Events

	if err := r.Db.Limit(10).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *Repository) GetBookingDataFromSignature(digitalSignature string) (*models.BookingData, error) {

	var bookingData models.BookingData

	if err := r.Db.First(&bookingData, "digital_signature =?", digitalSignature).Error; err != nil {
		return nil, err
	}
	return &bookingData, nil
}
