package helpers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// GetUserID — типизированный доступ к userID из контекста
// (используется во всех handler'ах)
func GetUserID(r *http.Request) (uuid.UUID, error) {
	val := r.Context().Value(UserIDKey)
	if val == nil {
		return uuid.Nil, fmt.Errorf("user id not found in context")
	}

	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user id type in context")
	}

	return id, nil
}
