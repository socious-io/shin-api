package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func Copy(src interface{}, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}

func DecodeJWT(token string) (header, payload []byte, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("invalid token format")
	}

	// Decode Header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding header: %v", err)
	}

	// Decode Payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding payload: %v", err)
	}

	return headerBytes, payloadBytes, nil
}

func ArrayContains[T comparable](arr []T, x T) bool {
	for _, item := range arr {
		if item == x {
			return true
		}
	}
	return false
}

func AppendIfNotExists[T comparable](arr []T, x T) []T {
	if !ArrayContains(arr, x) {
		arr = append(arr, x)
	}
	return arr
}
