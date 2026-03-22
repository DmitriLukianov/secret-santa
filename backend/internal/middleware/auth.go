package middleware

import (
	"context"
	"net/http"
	"strings"

	"secret-santa-backend/internal/auth"
)

type contextKey string

const userKey contextKey = "userID"

// Middleware проверки JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		}

		// Ожидаем формат: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// кладём userID в контекст
		ctx := context.WithValue(r.Context(), userKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Достаём userID из контекста
func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(userKey).(string)
	return userID, ok
}
