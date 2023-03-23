package encryption

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	edpublicKey  ed25519.PublicKey
	edprivateKey ed25519.PrivateKey
)

func GenerateEdDSAKey() {
	edpublicKey, edprivateKey, _ = ed25519.GenerateKey((nil))
}

func SignEdDsaMessage(message []byte) string {
	signature := ed25519.Sign(edprivateKey, message)
	sign := base64.StdEncoding.EncodeToString(signature)
	fmt.Printf("Signature: %v\n", sign)
	return string(sign)
}

func VerifyEdSignature(signature, message []byte) error {
	//digest := sha256.Sum256(message)
	res := ed25519.Verify(edpublicKey, message, signature)
	if !res {
		return errors.New("not able to verify")
	}
	return nil
}
