package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/spf13/viper"
	"github.com/yaserali542/ticket-management-service/encryption"
	"github.com/yaserali542/ticket-management-service/models"
	"github.com/yaserali542/ticket-management-service/repository"
)

type TicketManagementService struct {
	Repository *repository.Repository
}

func (*TicketManagementService) GetUserDetails(username, jwtToken string) (*models.BasicFields, error) {
	url := fmt.Sprintf("%v/basic-details", viper.GetViper().GetString("account_service_url"))

	requestBytes, _ := json.Marshal(models.BasicFieldsRequest{UserName: username})
	reqBody := bytes.NewReader(requestBytes)

	signature := encryption.CalculateHmac(requestBytes)

	req, _ := http.NewRequest(
		"POST",
		url,
		reqBody,
	)

	// add a request header
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Signature", signature)
	req.Header.Add("token", jwtToken)

	// send an HTTP using `req` object
	res, err := http.DefaultClient.Do(req)
	// check for response error
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid return status from :%v", url)
	}
	var basicAccountDetails models.BasicFields

	// read response data
	if err := json.NewDecoder(res.Body).Decode(&basicAccountDetails); err != nil {
		return nil, err
	}

	// close response body
	res.Body.Close()

	return &basicAccountDetails, nil

}

func (s *TicketManagementService) BookEvent(request models.BookEventRequest, userName string, jwtToken string) (*models.BookEventResponse, error) {

	userDetails, err := s.GetUserDetails(userName, jwtToken)

	if err != nil {
		return nil, err
	}

	event, err := s.Repository.GetEventDetails(request.EventId)
	if err != nil {
		return nil, err
	}
	bookEventResponse, _, err := s.Repository.GetBookingDetails(userDetails.ID, userName)

	if err != nil {
		return nil, err
	}

	if bookEventResponse != nil {
		for _, v := range bookEventResponse {
			if v.EventDetails.ID == request.EventId {
				return nil, errors.New("booking already exists for this user")
			}

		}
	}

	newBooking := &models.Booking{
		EventId: request.EventId,
		UserId:  userDetails.ID,
	}
	if err = s.Repository.CreateBooking(newBooking); err != nil {
		return nil, err
	}
	jsonBooking, _ := json.Marshal(newBooking)
	encryptedData := encryption.EncryptMessageUsingPublicKey(string(jsonBooking))
	hashEncryptedData := encryption.CalculateHash(encryptedData)
	//fmt.Println("hashedEncryptedData: ", string(hashEncryptedData))
	digitalSignature := encryption.SignMessage(hashEncryptedData)
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	bookingData := &models.BookingData{
		QrEncryptedData:  encryptedData,
		Otp:              otp,
		DigitalSignature: digitalSignature,
		BookingId:        newBooking.ID,
	}
	if err = s.Repository.CreateBookingData(bookingData); err != nil {
		return nil, err
	}

	return &models.BookEventResponse{
		EventDetails: models.Events{
			ID:         event.ID,
			EventImage: event.EventImage,
			StartDate:  event.StartDate,
			EventName:  event.EventName,
			EndDate:    event.EndDate,
			Venue:      event.Venue,
		},
		BookingId:  newBooking.ID,
		Signature:  digitalSignature,
		Otp:        otp,
		HashedData: base64.StdEncoding.EncodeToString(hashEncryptedData),
		UserId:     userDetails.ID.String(),
		UserName:   userDetails.UserName,
	}, nil
}

func (s *TicketManagementService) GetEvents() ([]models.Events, error) {
	return s.Repository.GetEvents()
}

func (s *TicketManagementService) GetBookingDetails(userName string, jwtToken string) ([]models.BookEventResponse, bool, error) {

	userDetails, err := s.GetUserDetails(userName, jwtToken)

	if err != nil {
		return nil, false, err
	}
	return s.Repository.GetBookingDetails(userDetails.ID, userName)
}

func (s *TicketManagementService) GetQRCodeDetails(digitalSignature, jwtToken string) (*models.QRCodeDetails, error) {

	bookingDetails, err := s.Repository.GetBookingDataFromSignature(digitalSignature)

	if err != nil {
		return nil, err
	}

	hashedData := encryption.CalculateHash(bookingDetails.QrEncryptedData)
	sig, _ := base64.StdEncoding.DecodeString(digitalSignature)
	if err := encryption.VerifySignature(sig, hashedData); err != nil {
		return nil, err
	}

	var booking models.Booking
	decryptedMessage := encryption.DecryptMessageUsingPrivateKey(bookingDetails.QrEncryptedData)

	json.Unmarshal([]byte(decryptedMessage), &booking)

	if booking.ID != bookingDetails.BookingId {
		return nil, errors.New("data mismatch: invalid booking and bookingDetails")
	}
	accountDetails, err := s.GetUserDetailsById(booking.UserId.String(), jwtToken)
	if err != nil {
		return nil, err
	}

	eventDetails, err := s.Repository.GetEventDetails(booking.EventId)
	if err != nil {
		return nil, err
	}

	return &models.QRCodeDetails{
		BookingDetails: booking,
		UserDetails:    *accountDetails,
		EventDetails:   *eventDetails,
		Otp:            bookingDetails.Otp,
	}, nil

}

func (*TicketManagementService) GetUserDetailsById(id, jwtToken string) (*models.Account, error) {
	url := fmt.Sprintf("%v/account/%v", viper.GetViper().GetString("account_service_url"), id)

	req, _ := http.NewRequest(
		"GET",
		url,
		nil,
	)

	// add a request header
	req.Header.Add("token", jwtToken)

	// send an HTTP using `req` object
	res, err := http.DefaultClient.Do(req)
	// check for response error
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid return status from :%v", url)
	}
	var accountDetails models.Account

	// read response data
	if err := json.NewDecoder(res.Body).Decode(&accountDetails); err != nil {
		return nil, err
	}

	// close response body
	res.Body.Close()

	return &accountDetails, nil

}
