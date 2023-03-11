package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

var (
	privateKey *rsa.PrivateKey = exportPEMStrToPrivKey(readKeyFromFile("privkey.pem"))
	publicKey  *rsa.PublicKey  = exportPEMStrToPubKey(readKeyFromFile("pubkey.pem"))
)

func readKeyFromFile(filename string) []byte {
	key, err := os.ReadFile(filename)

	if err != nil {
		fmt.Println("Error while reading the file :", filename, "Error :", err)
	}
	return key
}
func exportPEMStrToPrivKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return key
}
func exportPEMStrToPubKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	key, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return key
}

func EncryptMessageUsingPublicKey(message string) []byte {
	fmt.Println("message length to encrypt :,", len(message))
	cipherText, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		[]byte(message),
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("cipherText length", len(cipherText))
	return cipherText

}
func DecryptMessageUsingPrivateKey(cipherText []byte) string {
	decMessage, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipherText, nil)
	fmt.Printf("Original: %s\n", string(decMessage))
	if err != nil {
		panic(err)
	}

	fmt.Println("Encrypted message: ", string(cipherText))
	return string(cipherText)
}
func SignMessage(message []byte) string {
	//msgHashSum := CalculateHmac(message)
	// We have to provide a random reader, so every time
	// we sign, we have a different signature
	signature, _ := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, message, nil)
	sign := base64.StdEncoding.EncodeToString(signature)
	fmt.Printf("Signature: %v\n", sign)
	return string(sign)

}

func VerifySignature(signature, message []byte) error {
	msgHashSum := CalculateHmac(message)
	return rsa.VerifyPSS(publicKey, crypto.SHA256, []byte(msgHashSum), signature, nil)

}
