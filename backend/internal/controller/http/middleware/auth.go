package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"secret-santa-backend/internal/definitions"
	auth "secret-santa-backend/internal/oauth"

	"github.com/google/uuid"
)

type AuthMiddleware struct {
	jwtManager *auth.JWTManager
	log        *slog.Logger
}

func NewAuthMiddleware(jwtManager *auth.JWTManager, log *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager, log: log}
}

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtManager.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			http.Error(w, "invalid user id in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), definitions.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
