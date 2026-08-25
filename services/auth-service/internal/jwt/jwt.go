package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Manager is a thin wrapper around JWT signing using HMAC.
// In production you may want RS256 with rotating keys. HMAC is used here
// for simplicity in a learning project.
type Manager struct {
	secret []byte
}

// Claims is a minimal projection of the token payload used by the gRPC layer.
type Claims struct {
	UserID    string
	TokenType string
}

// New creates a JWT manager with a secret key.
func New(secret string) *Manager { return &Manager{secret: []byte(secret)} }

// NewAccessToken creates a short-lived access token for the given user ID.
func (m *Manager) NewAccessToken(userID string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"typ": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// NewRefreshToken creates a longer-lived refresh token for refreshing access tokens.
func (m *Manager) NewRefreshToken(userID string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"typ": "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Validate verifies a JWT and returns the parsed user identity used by internal RPCs.
func (m *Manager) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected token signature method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, _ := payload["sub"].(string)
	tokenType, _ := payload["typ"].(string)
	if userID == "" {
		return nil, errors.New("missing subject claim")
	}

	return &Claims{UserID: userID, TokenType: tokenType}, nil
}
