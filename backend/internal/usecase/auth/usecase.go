package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"log/slog"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/oauth"
	"secret-santa-backend/internal/usecase"
)

type UseCase struct {
	userUC           usecase.UserUseCase
	emailService     usecase.EmailService
	verificationRepo usecase.VerificationRepository
	smtpEnabled      bool // кэшируем статус при инициализации
	log              *slog.Logger
}

func New(userUC usecase.UserUseCase, emailService usecase.EmailService, verificationRepo usecase.VerificationRepository, smtpEnabled bool) *UseCase {
	return &UseCase{
		userUC:           userUC,
		emailService:     emailService,
		verificationRepo: verificationRepo,
		smtpEnabled:      smtpEnabled,
	}
}

func NewWithLogger(userUC usecase.UserUseCase, emailService usecase.EmailService, verificationRepo usecase.VerificationRepository, smtpEnabled bool, log *slog.Logger) *UseCase {
	return &UseCase{
		userUC:           userUC,
		emailService:     emailService,
		verificationRepo: verificationRepo,
		smtpEnabled:      smtpEnabled,
		log:              log,
	}
}

// LoginWithOAuth — вход через OAuth (GitHub и др.) + уведомление о входе
func (uc *UseCase) LoginWithOAuth(ctx context.Context, info oauth.UserInfo) (string, error) {
	if uc.log != nil {
		uc.log.Info("oauth login started",
			slog.String("provider", info.Provider),
			slog.String("oauth_id", info.ID),
		)
	}

	if info.ID == "" {
		return "", definitions.ErrMissingOAuthCode
	}

	user, err := uc.userUC.GetByOAuthID(ctx, info.ID, info.Provider)
	if err == nil && user != nil {
		if uc.log != nil {
			uc.log.Info("existing oauth user found", slog.String("user_id", user.ID.String()))
		}
		if uc.emailService != nil {
			_ = uc.emailService.SendLoginNotification(ctx, user.Email, user.Name)
		}
		return user.ID.String(), nil
	}

	if !errors.Is(err, definitions.ErrUserNotFound) && err != nil {
		return "", fmt.Errorf("failed to lookup oauth user: %w", err)
	}

	createInput := dto.CreateUserInput{
		Name:          info.Name,
		Email:         info.Email,
		OAuthID:       info.ID,
		OAuthProvider: info.Provider,
	}
	createdUser, err := uc.userUC.Create(ctx, createInput)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("new oauth user created", slog.String("user_id", createdUser.ID.String()))
	}
	if uc.emailService != nil {
		_ = uc.emailService.SendLoginNotification(ctx, createdUser.Email, createdUser.Name)
	}
	return createdUser.ID.String(), nil
}

// SendOTP — отправить код подтверждения на email.
// Явно возвращает ошибку если SMTP не настроен — OTP без email не имеет смысла.
func (uc *UseCase) SendOTP(ctx context.Context, email string) error {
	if uc.log != nil {
		uc.log.Info("send otp started", slog.String("email", email))
	}

	if !uc.smtpEnabled {
		return fmt.Errorf("email service is not configured: set SMTP_USERNAME, SMTP_PASSWORD and FROM_EMAIL")
	}

	code, err := uc.emailService.SendOTP(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to send OTP: %w", err)
	}

	expiresAt := time.Now().Add(10 * time.Minute)
	if err := uc.verificationRepo.SaveCode(ctx, email, code, expiresAt); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}

	return nil
}

// VerifyOTP — проверить код и вернуть userID.
func (uc *UseCase) VerifyOTP(ctx context.Context, email, code string) (string, error) {
	if uc.log != nil {
		uc.log.Info("verify otp started", slog.String("email", email))
	}

	valid, err := uc.verificationRepo.GetValidCode(ctx, email, code)
	if err != nil || !valid {
		return "", definitions.ErrInvalidUserInput
	}

	_ = uc.verificationRepo.MarkAsUsed(ctx, email, code)

	// Ищем существующего пользователя по email
	user, err := uc.userUC.GetByEmail(ctx, email)
	if err == nil && user != nil {
		if uc.emailService != nil {
			_ = uc.emailService.SendLoginNotification(ctx, user.Email, user.Name)
		}
		return user.ID.String(), nil
	}

	// Если пользователя нет — создаём (passwordless регистрация)
	// Берём имя из email (часть до @) как дефолтное
	defaultName := nameFromEmail(email)

	createInput := dto.CreateUserInput{
		Name:          defaultName,
		Email:         email,
		OAuthID:       email,
		OAuthProvider: "email",
	}
	createdUser, err := uc.userUC.Create(ctx, createInput)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	if uc.emailService != nil {
		_ = uc.emailService.SendLoginNotification(ctx, createdUser.Email, createdUser.Name)
	}
	return createdUser.ID.String(), nil
}

// nameFromEmail — извлекает часть до @ и делает из неё читаемое имя.
// "john.doe@example.com" → "john.doe"
func nameFromEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) == 0 || parts[0] == "" {
		return "Пользователь"
	}
	return parts[0]
}
