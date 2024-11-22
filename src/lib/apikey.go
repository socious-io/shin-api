package lib

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func GenerateSecureKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateApiKey() string {
	const apiKeyLength = 32

	apiKey, err := GenerateSecureKey(apiKeyLength)
	if err != nil {
		log.Fatalf("Failed to generate API key: %v", err)
	}

	return apiKey
}

func GenerateApiSecret() string {
	const apiSecretLength = 64

	apiSecret, err := GenerateSecureKey(apiSecretLength)
	if err != nil {
		log.Fatalf("Failed to generate API secret: %v", err)
	}

	return apiSecret
}
