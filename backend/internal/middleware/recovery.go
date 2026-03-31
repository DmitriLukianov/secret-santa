package middleware

import (
	"fmt"
	"net/http"

	"secret-santa-backend/internal/controller/http/v1/response"
)

// RecoveryMiddleware — защищает от panic
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err := fmt.Errorf("panic: %v", rec)
				response.WriteHTTPError(w, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
