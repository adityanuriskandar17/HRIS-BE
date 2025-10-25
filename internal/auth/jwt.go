package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func SignJWT(uid uint64, role string, secret string, exp time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  fmt.Sprintf("%d", uid),
		"role": role,
		"exp":  time.Now().Add(exp).Unix(),
		"iat":  time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}
