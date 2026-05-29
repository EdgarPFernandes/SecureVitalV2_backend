package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(
			tokenString,
			&Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			},
		)

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*Claims)

		ctx := context.WithValue(
			r.Context(),
			UserIDKey,
			claims.UserID,
		)

		ctx = context.WithValue(
			ctx,
			RoleKey,
			claims.Role,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
