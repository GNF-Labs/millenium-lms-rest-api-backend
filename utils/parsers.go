package utils

import (
	"errors"
	"strings"
)

func ParseToken(bearerToken string) (string, error) {
	if bearerToken == "" {
		return "", errors.New("token is empty")
	}

	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) != 2 {
		//c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return "", errors.New("token format error")
	}

	tokenString := splitToken[1]
	return tokenString, nil
}
