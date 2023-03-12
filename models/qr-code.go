package models

type VerifyQRCodeRequest struct {
	DigitalSignature string `json:"digital_signature"`
}

type QRCodeDetails struct {
	UserDetails    Account `json:"user_details"`
	EventDetails   Events  `json:"event_details"`
	BookingDetails Booking `json:"booking_details"`
	Otp            string  `json:"otp"`
}
