package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/yaserali542/ticket-management-service/models"
	"github.com/yaserali542/ticket-management-service/services"
)

type Middleware struct {
	Services services.TicketManagementService
}

func (m *Middleware) ValidateRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ServeHTTP(w, r)
			return
		}
		secret := viper.GetViper().GetString("hmac_secret")
		//secret := "mysecret"

		byteData, err := io.ReadAll(r.Body)

		if err != nil {
			errMsg := "Bad data"
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
		r.Body.Close() //  must close
		r.Body = io.NopCloser(bytes.NewBuffer(byteData))
		hmac := hmac.New(sha256.New, []byte(secret))

		hmac.Write(byteData)

		sha := hex.EncodeToString(hmac.Sum(nil))

		sig := r.Header.Get("Signature")

		if sha != sig {
			errMsg := "Forbidden"
			http.Error(w, errMsg, http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)

	})
}
func (m *Middleware) ValidateJwtToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := viper.GetViper()

		//url := fmt.Sprintf("%v/basic-details", v.GetString("account_service_url"))

		jwtToken := r.Header.Get("token")

		if jwtToken == "" {
			http.Error(w, "token is missing", http.StatusBadRequest)
			return
		}
		claims := &models.Claims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
			jwtKey := v.GetString("jwt.key")
			return []byte(jwtKey), nil
		})
		fmt.Println(err)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !(time.Until(claims.ExpiresAt.Time) > time.Duration(v.GetInt("jwt.expire-time-minutes")-v.GetInt("jwt.refresh-time-minutes"))*time.Minute) {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		var basicAccountDetails *models.BasicFields
		basicAccountDetails, err = m.Services.GetUserDetails(claims.Username, jwtToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if basicAccountDetails.UserName != claims.Username {
			http.Error(w, "username not found", http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(w, r)

	})

}
