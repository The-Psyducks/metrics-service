package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

var (
	jwtSecret          string
	jwtExpirationHours int
)

func init() {
	jwtSecret = os.Getenv("JWT_SECRET")
	jwtExpirationHoursStr := os.Getenv("JWT_DURATION_HOURS")

	if jwtExpirationHoursStr != "" {
		jwtExpirationHours, _ = strconv.Atoi(jwtExpirationHoursStr)
	}
}

type Claims struct {
	UserId    string `json:"user_id"`
	UserAdmin bool   `json:"user_admin"`
	jwt.RegisteredClaims
}

func GenerateToken(userId string, userAdmin bool) (string, error) {
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable is not set")
	}
	if jwtExpirationHours == 0 {
		return "", fmt.Errorf("JWT_DURATION_HOURS environment variable is not set or invalid")
	}

	claims := Claims{
		UserId:    userId,
		UserAdmin: userAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ThePsyducks-users-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sign, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("error signing JWT token: %w", err)
	}

	return sign, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing JWT token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}
