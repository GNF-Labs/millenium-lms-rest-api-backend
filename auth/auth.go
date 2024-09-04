// Package auth provides functions for generating and verifying JWT tokens for user authentication.
package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Claims represents the JWT claims, which includes the username and registered claims.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JSON Web Token (JWT) for a user after login.
//
// The token is signed using the provided jwtKey and includes the username
// in the claims. The token expires 24 hours after it is issued.
//
// Parameters:
//   - jwtKey: A byte slice containing the key to sign the JWT.
//   - username: The username to include in the JWT claims.
//
// Returns:
//   - A string representing the signed JWT.
//   - An error if there is an issue generating the token.
func GenerateJWT(jwtKey []byte, username string) (string, error) {
	expirationTime := time.Now().Add(24 * 7 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

// VerifyJWT verifies a JSON Web Token (JWT) and returns the claims if the token is valid.
//
// The token is parsed and validated using the provided jwtKey. If the token is valid
// and the signature is correct, the claims are returned. If the token is invalid or
// the signature does not match, an error is returned.
//
// Parameters:
//   - jwtKey: A byte slice containing the key to verify the JWT signature.
//   - tokenString: The JWT string to be verified.
//
// Returns:
//   - A pointer to the Claims struct if the token is valid.
//   - An error if the token is invalid or the signature does not match.
func VerifyJWT(jwtKey []byte, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
