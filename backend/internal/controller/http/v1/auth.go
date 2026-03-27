package v1

import (
	"encoding/json"
	"net/http"
	"time"

	authpkg "secret-santa-backend/internal/auth"
	"secret-santa-backend/internal/definitions"
	authuc "secret-santa-backend/internal/usecase/auth"
)

type AuthHandler struct {
	provider authpkg.Provider
	jwt      *authpkg.JWTManager
	uc       *authuc.UseCase
}

func NewAuthHandler(provider authpkg.Provider, jwt *authpkg.JWTManager, uc *authuc.UseCase) *AuthHandler {
	return &AuthHandler{
		provider: provider,
		jwt:      jwt,
		uc:       uc,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state := "ss-" + time.Now().Format("20060102150405")

	url := h.provider.Config().AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.URL.Query().Get("code")
	if code == "" {
		writeHTTPError(w, definitions.ErrInvalidOAuthCode)
		return
	}

	token, err := h.provider.Config().Exchange(ctx, code)
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	user, err := h.provider.GetUserInfo(ctx, token)
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	userID, err := h.uc.LoginWithOAuth(ctx, user)
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	jwtToken, err := h.jwt.GenerateToken(userID)
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": jwtToken})
}
