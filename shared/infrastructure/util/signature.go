package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func CreateSignature(clientSecret string, payloadInBytes []byte) (string, error) {
	h := hmac.New(sha256.New, []byte(clientSecret))
	_, err := h.Write(payloadInBytes)
	if err != nil {
		return "", err
	}
	signature := fmt.Sprintf("%x", h.Sum(nil))
	return signature, nil
}

func CreateSHA256Signature(stringSignature string) string {
	sum := sha256.Sum256([]byte(stringSignature))
	return fmt.Sprintf("%x", sum)
}
