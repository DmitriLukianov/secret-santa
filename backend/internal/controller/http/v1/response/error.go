package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"secret-santa-backend/internal/definitions"
)

// ErrorResponse — единый формат ошибки для всего API
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// WriteHTTPError — главная функция обработки ошибок (экспортирована)
func WriteHTTPError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	status := http.StatusInternalServerError
	message := err.Error()

	switch {
	case errors.Is(err, definitions.ErrNotFound),
		errors.Is(err, definitions.ErrEventNotFound),
		errors.Is(err, definitions.ErrWishlistNotFound),
		errors.Is(err, definitions.ErrParticipantNotFound),
		errors.Is(err, definitions.ErrAssignmentNotFound),
		errors.Is(err, definitions.ErrUserNotFound):
		status = http.StatusNotFound

	case errors.Is(err, definitions.ErrForbidden),
		errors.Is(err, definitions.ErrNotOrganizer),
		errors.Is(err, definitions.ErrNotSanta),
		errors.Is(err, definitions.ErrWishlistVisibilityForbidden):
		status = http.StatusForbidden

	case errors.Is(err, definitions.ErrConflict),
		errors.Is(err, definitions.ErrAlreadyParticipating),
		errors.Is(err, definitions.ErrDuplicateParticipant),
		errors.Is(err, definitions.ErrEventAlreadyFinished):
		status = http.StatusConflict

	case errors.Is(err, definitions.ErrInvalidUUID),
		errors.Is(err, definitions.ErrInvalidOAuthCode),
		errors.Is(err, definitions.ErrMissingOAuthCode),
		errors.Is(err, definitions.ErrInvalidUserInput),
		errors.Is(err, definitions.ErrInvalidOAuthUserInfo),
		errors.Is(err, definitions.ErrInvalidEventState),
		errors.Is(err, definitions.ErrNotEnoughParticipants),
		errors.Is(err, definitions.ErrInvalidWishlistVisibility):
		status = http.StatusBadRequest

	default:
		// Для неизвестных ошибок не раскрываем детали
		status = http.StatusInternalServerError
		message = "internal server error"
	}

	writeJSONError(w, status, message)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  status,
	})
}
