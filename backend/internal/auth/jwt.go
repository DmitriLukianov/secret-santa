package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager управляет генерацией и парсингом JWT-токенов.
type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

// NewJWTManager создаёт JWTManager. Падает с ошибкой, если секрет пуст или слишком короткий.
func NewJWTManager(secret string, ttl time.Duration) (*JWTManager, error) {
	if secret == "" {
		return nil, errors.New("jwt secret must not be empty")
	}
	if len(secret) < 32 {
		return nil, errors.New("jwt secret must be at least 32 characters")
	}
	return &JWTManager{secret: []byte(secret), ttl: ttl}, nil
}

// Claims — payload JWT-токена.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken создаёт подписанный JWT-токен для пользователя.
func (m *JWTManager) GenerateToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "secret-santa",
			Audience:  jwt.ClaimStrings{"secret-santa-api"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken проверяет и парсит JWT-токен.
func (m *JWTManager) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return m.secret, nil
	},
		jwt.WithIssuer("secret-santa"),
		jwt.WithAudience("secret-santa-api"),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
