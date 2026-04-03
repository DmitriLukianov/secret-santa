package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort  string `env:"APP_PORT" envDefault:"8080"`
	AppEnv   string `env:"APP_ENV" envDefault:"local"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	DatabaseURL string `env:"DATABASE_URL"`

	JWTSecret string        `env:"JWT_SECRET"`
	JWTTTL    time.Duration `env:"JWT_TTL" envDefault:"24h"`

	// === SMTP Mail.ru ===
	SMTPHost     string `env:"SMTP_HOST" envDefault:"smtp.mail.ru"`
	SMTPPort     int    `env:"SMTP_PORT" envDefault:"587"`
	SMTPUsername string `env:"SMTP_USERNAME"` // ваш полный email, например: vasya@mail.ru
	SMTPPassword string `env:"SMTP_PASSWORD"` // приложенный пароль (НЕ основной!)
	FromEmail    string `env:"FROM_EMAIL"`    // от кого будут письма (обычно тот же email)

	OAuthProvider      string `env:"OAUTH_PROVIDER" envDefault:"github"`
	GithubClientID     string `env:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	GithubRedirectURL  string `env:"GITHUB_REDIRECT_URL"`
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		AppEnv:      getEnv("APP_ENV", "local"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		JWTTTL:      parseDuration(getEnv("JWT_TTL", "24h")),

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

	// Обязательные проверки
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" || len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters")
	}
	if cfg.SMTPUsername == "" || cfg.SMTPPassword == "" || cfg.FromEmail == "" {
		log.Fatal("SMTP_USERNAME, SMTP_PASSWORD and FROM_EMAIL are required for email notifications")
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
