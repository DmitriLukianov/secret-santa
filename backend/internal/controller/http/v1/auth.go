package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/helpers"
	"secret-santa-backend/internal/oauth"
	"secret-santa-backend/internal/usecase"
	authuc "secret-santa-backend/internal/usecase/auth"
)

type AuthHandler struct {
	provider    oauth.Provider
	jwt         *oauth.JWTManager
	uc          *authuc.UseCase
	userUC      usecase.UserUseCase
	frontendURL string
}

func NewAuthHandler(provider oauth.Provider, jwt *oauth.JWTManager, uc *authuc.UseCase, userUC usecase.UserUseCase, frontendURL string) *AuthHandler {
	return &AuthHandler{
		provider:    provider,
		jwt:         jwt,
		uc:          uc,
		userUC:      userUC,
		frontendURL: frontendURL,
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

	redirectURL := h.frontendURL + "/auth/callback?token=" + jwtToken
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *AuthHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.SendOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	if err := h.uc.SendOTP(r.Context(), req.Email); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Код отправлен на почту",
	})
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	userID, err := h.uc.VerifyOTP(r.Context(), req.Email, req.Code, req.Name)
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

	parsedID, err := uuid.Parse(userID)
	if err == nil {
		user, err := h.userUC.GetByID(r.Context(), parsedID)
		if err == nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"token": jwtToken,
				"user":  response.UserResponse{ID: user.ID.String(), Name: user.Name, Email: user.Email},
			})
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string{"token": jwtToken})
}
