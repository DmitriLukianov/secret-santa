package helpers

import "strings"

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	return strings.Contains(errStr, "duplicate key value violates unique constraint") ||
		strings.Contains(errStr, "SQLSTATE 23505")
}
