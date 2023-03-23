package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/yaserali542/ticket-management-service/models"
	"github.com/yaserali542/ticket-management-service/services"
)

type Controllers struct {
	Services services.TicketManagementService
}

func (c *Controllers) BookEvent(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UnixNano() / int64(time.Millisecond)

	// do something

	var request models.BookEventRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "incorrect data", http.StatusBadRequest)
		return
	}
	claims := &models.Claims{}
	jwtToken := r.Header.Get("token")

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		jwtKey := viper.GetViper().GetString("jwt.key")
		return []byte(jwtKey), nil
	})

	response, err := c.Services.BookEvent(request, claims.Username, jwtToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
	end := time.Now().UnixNano() / int64(time.Millisecond)
	diff := end - start
	fmt.Printf("Duration(ms): %d", diff)
}

func (c *Controllers) GetEvents(w http.ResponseWriter, r *http.Request) {

	events, err := c.Services.GetEvents()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(events)
}

func (c *Controllers) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	claims := &models.Claims{}
	jwtToken := r.Header.Get("token")

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		jwtKey := viper.GetViper().GetString("jwt.key")
		return []byte(jwtKey), nil
	})

	response, errRecordNotFound, err := c.Services.GetBookingDetails(claims.Username, jwtToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if errRecordNotFound {
		http.Error(w, "no records found", http.StatusNoContent)
	}

	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (c *Controllers) GetBookingFromId(w http.ResponseWriter, r *http.Request) {
	claims := &models.Claims{}
	jwtToken := r.Header.Get("token")
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		jwtKey := viper.GetViper().GetString("jwt.key")
		return []byte(jwtKey), nil
	})

	response, errRecordNotFound, err := c.Services.GetBookingDetails(claims.Username, jwtToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if errRecordNotFound {
		http.Error(w, "no records found", http.StatusNoContent)
	}

	for _, v := range response {
		if v.BookingId.String() == id {
			w.Header().Add("Content-type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(v)
			return
		}
	}
	if errRecordNotFound {
		http.Error(w, "no record with id found", http.StatusNoContent)
	}
}

// GetQRCodeDetails
func (c *Controllers) GetQRCodeDetails(w http.ResponseWriter, r *http.Request) {

	var request models.VerifyQRCodeRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "incorrect data", http.StatusBadRequest)
		return
	}
	jwtToken := r.Header.Get("token")

	response, err := c.Services.GetQRCodeDetails(request.DigitalSignature, jwtToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)

}
