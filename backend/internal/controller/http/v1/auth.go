package v1

import (
	"net/http"

	oauth "secret-santa-backend/internal/auth" // ← JWT + provider
	"secret-santa-backend/internal/usecase/auth"
)

type AuthHandler struct {
	provider oauth.Provider // 🔥 ВОТ ЭТО ГЛАВНОЕ ИЗМЕНЕНИЕ
	uc       *auth.UseCase
}

func NewAuthHandler(provider oauth.Provider, uc *auth.UseCase) *AuthHandler {
	return &AuthHandler{
		provider: provider,
		uc:       uc,
	}
}

// GET /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	url := h.provider.Config().AuthCodeURL("state") // 🔥 небольшое упрощение
	http.Redirect(w, r, url, http.StatusFound)
}

// GET /auth/callback
func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code not found", http.StatusBadRequest)
		return
	}

	token, err := h.provider.Config().Exchange(ctx, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.provider.GetUserInfo(ctx, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 👇 теперь получаем userID
	userID, err := h.uc.LoginWithOAuth(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 👇 создаём JWT
	jwtToken, err := oauth.GenerateToken(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 👇 возвращаем токен
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"token":"` + jwtToken + `"}`))
}
