package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key")

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

func GenerateToken(userID int64, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(jwtSecret)
}
