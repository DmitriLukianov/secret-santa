package email

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"net/smtp"

	"secret-santa-backend/internal/config"
)

type Service struct {
	cfg *config.Config
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *Service {
	return &Service{cfg: cfg, log: log}
}

// SendLoginNotification — уведомление при входе через GitHub
func (s *Service) SendLoginNotification(ctx context.Context, email, name string) error {
	subject := "✅ Вы успешно вошли в Тайный Санта"
	body := fmt.Sprintf(`Привет, %s!

Вы только что вошли в аккаунт Тайный Санта.

Если это были не вы — немедленно напишите нам.

С наилучшими пожеланиями,
Команда Тайный Санта`, name)

	return s.send(ctx, email, subject, body)
}

// SendOTP — отправка 6-значного кода
func (s *Service) SendOTP(ctx context.Context, email string) (string, error) {
	code := s.generateOTP()

	subject := "🔑 Код подтверждения для Тайный Санта"
	body := fmt.Sprintf(`Ваш код подтверждения: %s

Код действителен 10 минут.

Если вы не запрашивали код — просто проигнорируйте это письмо.

С наилучшими пожеланиями,
Команда Тайный Санта`, code)

	err := s.send(ctx, email, subject, body)
	if err != nil {
		return "", err
	}
	return code, nil
}

// === НОВОЕ: Уведомление при начале жеребьёвки ===
func (s *Service) SendDrawNotification(ctx context.Context, email, eventTitle string) error {
	subject := "🎲 Жеребьёвка проведена!"
	body := fmt.Sprintf(`Привет!

В событии «%s» проведена жеребьёвка.

Теперь вы можете посмотреть, кому дарите подарок, в разделе "Мои события".

Удачи в Тайном Санте! 🎁

С наилучшими пожеланиями,
Команда Тайный Санта`, eventTitle)

	return s.send(ctx, email, subject, body)
}

// generateOTP — 6-значный код
func (s *Service) generateOTP() string {
	const digits = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code[i] = digits[n.Int64()]
	}
	return string(code)
}

func (s *Service) send(ctx context.Context, to, subject, body string) error {
	if s.log != nil {
		s.log.Info("sending email", slog.String("to", to), slog.String("subject", subject))
	}

	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n"+
		"%s\r\n", s.cfg.FromEmail, to, subject, body))

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	if err := smtp.SendMail(addr, auth, s.cfg.FromEmail, []string{to}, msg); err != nil {
		if s.log != nil {
			s.log.Error("failed to send email", slog.String("to", to), slog.String("error", err.Error()))
		}
		return fmt.Errorf("failed to send email: %w", err)
	}

	if s.log != nil {
		s.log.Info("email sent successfully", slog.String("to", to))
	}
	return nil
}
