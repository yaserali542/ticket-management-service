package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/yaserali542/ticket-management-service/models"
	"github.com/yaserali542/ticket-management-service/services"
)

type Controllers struct {
	Services services.TicketManagementService
}

func (c *Controllers) BookEvent(w http.ResponseWriter, r *http.Request) {

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
}
