package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	// "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v4"
	// "github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "cineverse_secret_key"
	}
	jwtSecret = []byte(secret)
}

// Claims struct
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// Access token (1 hour)
func CreateAccessToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func CreateRefreshToken(userID uint) (string, error) {
	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		return "", errors.New("missing REFRESH_SECRET environment variable")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshSecret))
}

// func ParseToken(tokenStr string) (*Claims, error) {
// 	tkn, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return jwtSecret, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if claims, ok := tkn.Claims.(*Claims); ok && tkn.Valid {
// 		return claims, nil
// 	}
// 	return nil, errors.New("invalid token")
// }

// ValidateRefreshToken validates the refresh token and returns claims if valid
func ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		return nil, errors.New("missing REFRESH_SECRET environment variable")
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Verify the signing method is HMAC
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(refreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration manually
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return nil, errors.New("refresh token expired")
			}
		}
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}

// func ValidateAccessToken(tokenStr string) (*Claims, error) {
// 	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
// 		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
// 		}
// 		return jwtSecret, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		return claims, nil
// 	}
// 	return nil, errors.New("invalid access token")
// }

// ValidateAccessToken validates the access token and returns its claims if valid.
func ValidateAccessToken(tokenStr string) (*Claims, error) {
	// Parse the token with the custom Claims struct
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// Explicitly check for HS256 signing method
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v, expected HS256", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		// Handle specific JWT validation errors
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("malformed access token")
			}
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, errors.New("access token is expired or not yet valid")
			}
			return nil, fmt.Errorf("access token validation failed: %w", err)
		}
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	// Check if token is valid and claims are of the expected type
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Additional validation of critical claims
		if claims.UserID == 0 {
			return nil, errors.New("invalid access token: missing user_id claim")
		}
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("access token expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid access token: claims are invalid or token is not valid")
}

// func ValidateAccessToken(tokenStr string, jwtSecret []byte) (*Claims, error) {
// 	// Parse the token with the custom Claims struct
// 	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
// 		// Explicitly check for HS256 signing method
// 		if t.Method != jwt.SigningMethodHS256 {
// 			return nil, fmt.Errorf("unexpected signing method: %v, expected HS256", t.Header["alg"])
// 		}
// 		return jwtSecret, nil
// 	})

// 	if err != nil {
// 		// Handle specific JWT validation errors
// 		if ve, ok := err.(*jwt.ValidationError); ok {
// 			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
// 				return nil, errors.New("malformed token")
// 			}
// 			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
// 				return nil, errors.New("token is expired or not yet valid")
// 			}
// 			return nil, fmt.Errorf("token validation failed: %w", err)
// 		}
// 		return nil, fmt.Errorf("failed to parse token: %w", err)
// 	}

// 	// Check if token is valid and claims are of the expected type
// 	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
// 		// Additional explicit validation of critical claims
// 		if claims. != 0 && claims.ExpiresAt < jwt.TimeFunc().Unix() {
// 			return nil, errors.New("token is expired")
// 		}
// 		if claims.UserID == "" {
// 			return nil, errors.New("invalid token: missing user_id claim")
// 		}
// 		return claims, nil
// 	}

// 	return nil, errors.New("invalid token: claims are invalid or token is not valid")
// }
