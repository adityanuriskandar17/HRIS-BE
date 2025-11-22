package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	Role     string `json:"role"`
	TenantID string `json:"tenant_id,omitempty"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

func SignAccessToken(uid uint64, role string, secret string, ttl time.Duration, tenantID string) (string, time.Time, error) {
	now := time.Now().UTC()
	claims := AccessTokenClaims{
		Role:     role,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", uid),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, claims.ExpiresAt.Time, nil
}

func ParseAccessToken(token string, secret string) (*AccessTokenClaims, error) {
	t, err := jwt.ParseWithClaims(token, &AccessTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %q", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(*AccessTokenClaims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// Refresh token helpers ----------------------------------------------------

const refreshTokenLength = 32

// GenerateRefreshToken produces a random refresh token and returns the plaintext value,
// its hash (for storage), and the expiry timestamp.
func GenerateRefreshToken(ttl time.Duration) (plain string, hash string, expiresAt time.Time, err error) {
	if ttl <= 0 {
		return "", "", time.Time{}, errors.New("refresh ttl must be > 0")
	}
	buf := make([]byte, refreshTokenLength)
	if _, err = rand.Read(buf); err != nil {
		return "", "", time.Time{}, err
	}
	plain = base64.RawURLEncoding.EncodeToString(buf)
	sum := sha256.Sum256([]byte(plain))
	expiresAt = time.Now().UTC().Add(ttl)
	return plain, fmt.Sprintf("%x", sum[:]), expiresAt, nil
}

// HashRefreshToken converts a plaintext refresh token into a sha256 hash for lookup.
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", sum[:])
}
