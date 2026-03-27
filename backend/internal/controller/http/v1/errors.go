package v1

import (
	"errors"
	"net/http"

	"secret-santa-backend/internal/definitions"
)

func writeHTTPError(w http.ResponseWriter, err error) {
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
	}

	http.Error(w, message, status)
}
