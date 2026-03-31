package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/oauth"
	authuc "secret-santa-backend/internal/usecase/auth"
)

type AuthHandler struct {
	provider oauth.Provider
	jwt      *oauth.JWTManager
	uc       *authuc.UseCase
}

func NewAuthHandler(provider oauth.Provider, jwt *oauth.JWTManager, uc *authuc.UseCase) *AuthHandler {
	return &AuthHandler{
		provider: provider,
		jwt:      jwt,
		uc:       uc,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state := "ss-" + time.Now().Format("20060102150405")
	url := h.provider.GetAuthURL(state)

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.URL.Query().Get("code")
	if code == "" {
		response.WriteHTTPError(w, definitions.ErrInvalidOAuthCode)
		return
	}

	token, err := h.provider.Config().Exchange(ctx, code)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	userInfo, err := h.provider.GetUserInfo(ctx, token)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	userID, err := h.uc.LoginWithOAuth(ctx, userInfo)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	jwtToken, err := h.jwt.GenerateToken(userID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": jwtToken})
}
