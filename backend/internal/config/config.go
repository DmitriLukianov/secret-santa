package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string `env:"APP_PORT" envDefault:"8080"`
	AppEnv     string `env:"APP_ENV" envDefault:"local"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"info"`
	AppBaseURL  string `env:"APP_BASE_URL" envDefault:"http://localhost:8080"`
	FrontendURL string `env:"FRONTEND_URL" envDefault:"http://localhost:5173"`

	DatabaseURL string `env:"DATABASE_URL"`

	JWTSecret string        `env:"JWT_SECRET"`
	JWTTTL    time.Duration `env:"JWT_TTL" envDefault:"24h"`

	// SMTP — опциональный. Если не задан, email-уведомления отключаются.
	// Сервис стартует без них. Проверяй через cfg.SMTPEnabled().
	SMTPHost     string `env:"SMTP_HOST" envDefault:"smtp.mail.ru"`
	SMTPPort     int    `env:"SMTP_PORT" envDefault:"587"`
	SMTPUsername string `env:"SMTP_USERNAME"`
	SMTPPassword string `env:"SMTP_PASSWORD"`
	FromEmail    string `env:"FROM_EMAIL"`

	OAuthProvider      string `env:"OAUTH_PROVIDER" envDefault:"github"`
	GithubClientID     string `env:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	GithubRedirectURL  string `env:"GITHUB_REDIRECT_URL"`
}

// SMTPEnabled возвращает true если все необходимые SMTP-переменные заданы.
// Используй это перед отправкой писем вместо падения сервиса.
func (c *Config) SMTPEnabled() bool {
	return c.SMTPUsername != "" && c.SMTPPassword != "" && c.FromEmail != ""
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "local"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		AppBaseURL:  getEnv("APP_BASE_URL", "http://localhost:8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),

		DatabaseURL: getEnv("DATABASE_URL", ""),

		JWTSecret: getEnv("JWT_SECRET", ""),
		JWTTTL:    parseDuration(getEnv("JWT_TTL", "24h")),

		SMTPHost:     getEnv("SMTP_HOST", "smtp.mail.ru"),
		SMTPPort:     getIntEnv("SMTP_PORT", 587),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", ""),

		OAuthProvider:      getEnv("OAUTH_PROVIDER", "github"),
		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GithubRedirectURL:  getEnv("GITHUB_REDIRECT_URL", ""),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" || len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters")
	}

	if !cfg.SMTPEnabled() {
		log.Println("WARNING: SMTP not configured — email notifications are disabled")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("Warning: invalid JWT_TTL '%s', using default 24h", s)
		return 24 * time.Hour
	}
	return d
}
