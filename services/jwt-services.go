package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/yaserali542/ticket-management-service/models"
)

func GenerateToken(username string) string {
	v := viper.GetViper()
	//expirationTime := time.Now().Add(v.GetInt("jwt.key") * time.Minute)
	expirationTime := time.Now().Add(time.Duration(v.GetInt("jwt.expire-time-minutes")) * time.Minute)

	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	jwtSecret := v.GetString("jwt.key")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, _ := token.SignedString([]byte(jwtSecret))

	return tokenString
}
