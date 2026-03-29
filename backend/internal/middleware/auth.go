package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	auth "secret-santa-backend/internal/oauth"

	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "user_id"

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
			if m.log != nil {
				m.log.Warn("missing authorization header")
			}
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			if m.log != nil {
				m.log.Warn("invalid authorization header format")
			}
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		if m.log != nil {
			m.log.Info("parsing token", slog.String("token_prefix", tokenStr[:15]+"..."))
		}

		claims, err := m.jwtManager.ParseToken(tokenStr)
		if err != nil {
			if m.log != nil {
				m.log.Error("failed to parse token", slog.String("error", err.Error()))
			}
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if claims.UserID == "" {
			if m.log != nil {
				m.log.Error("token has no user_id")
			}
			http.Error(w, "invalid user id in token", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			if m.log != nil {
				m.log.Error("failed to parse user_id", slog.String("raw_user_id", claims.UserID))
			}
			http.Error(w, "invalid user id in token", http.StatusUnauthorized)
			return
		}

		if m.log != nil {
			m.log.Info("✅ user authenticated", slog.String("user_id", userID.String()))
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) (uuid.UUID, error) {
	val := r.Context().Value(UserIDKey)
	if val == nil {
		return uuid.Nil, fmt.Errorf("user id not found in context")
	}
	return val.(uuid.UUID), nil
}
