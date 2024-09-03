package utils

import (
	"encoding/base64"
	"errors"
	"strings"
)

func DecodeBase64ToBytes(encoded string) ([]byte, error) {
	if encoded == "" {
		return nil, errors.New("input base64 string is empty")
	}

	if strings.HasPrefix(encoded, "data:") {
		// Split the data URL into metadata and the base64 data
		parts := strings.SplitN(encoded, ",", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid data URL format")
		}
		// The actual base64 encoded data is in the second part
		encoded = parts[1]
	}

	// Decode the base64 string
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	return decodedBytes, nil
}
