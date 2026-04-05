package middleware

import (
	"context"
	"encoding/json"
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
			writeJSONUnauthorized(w, "missing authorization header")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			writeJSONUnauthorized(w, "invalid authorization header format")
			return
		}

		claims, err := m.jwtManager.ParseToken(tokenStr)
		if err != nil {
			writeJSONUnauthorized(w, "invalid token")
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			writeJSONUnauthorized(w, "invalid user id in token")
			return
		}

		ctx := context.WithValue(r.Context(), definitions.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeJSONUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": msg,
		"code":  http.StatusUnauthorized,
	})
}
