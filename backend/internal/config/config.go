package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Приложение
	AppPort  string `env:"APP_PORT" envDefault:"8080"`
	AppEnv   string `env:"APP_ENV" envDefault:"local"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// База данных
	DatabaseURL string `env:"DATABASE_URL"`

	// JWT
	JWTSecret string        `env:"JWT_SECRET"`
	JWTTTL    time.Duration `env:"JWT_TTL" envDefault:"24h"`

	// OAuth (поддержка нескольких провайдеров)
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
		AppPort:            getEnv("APP_PORT", "8080"),
		AppEnv:             getEnv("APP_ENV", "local"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTTTL:             parseDuration(getEnv("JWT_TTL", "24h")),
		OAuthProvider:      getEnv("OAUTH_PROVIDER", "github"),
		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GithubRedirectURL:  getEnv("GITHUB_REDIRECT_URL", ""),
	}

	// Простая валидация обязательных полей
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required and must be at least 32 characters")
	}
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
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
