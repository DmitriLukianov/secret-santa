package middleware

import (
	"context"
	"net/http"
	"strings"

	"secret-santa-backend/internal/auth"
)

type contextKey string

const userKey contextKey = "userID"

// AuthMiddleware — структура для middleware проверки JWT.
type AuthMiddleware struct {
	jwt *auth.JWTManager
}

// NewAuthMiddleware создаёт middleware с инжектированным JWTManager.
func NewAuthMiddleware(jwt *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwt}
}

// Handler проверяет JWT-токен и кладёт userID в контекст.
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.jwt.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID достаёт userID из контекста.
func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(userKey).(string)
	return userID, ok
}
