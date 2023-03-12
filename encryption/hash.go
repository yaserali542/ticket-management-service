package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/spf13/viper"
)

func CalculateHmac(data []byte) string {

	secret := viper.GetViper().GetString("hmac_secret")
	hmac := hmac.New(sha256.New, []byte(secret))

	hmac.Write(data)
	hmacData := hex.EncodeToString(hmac.Sum(nil))
	return hmacData
}

func CalculateHash(data []byte) []byte {
	msgHash := sha256.New()
	msgHash.Write(data)
	msgHashSum := msgHash.Sum(nil)
	return msgHashSum
}