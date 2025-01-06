package jwt

import (
	"errors"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	tokenSecret   = []byte("g0KKCTHGBeX8YDLGq3eOHZDvVBLuhz3p")
	refreshSecret = []byte("RySA4AXIMYr4bJEPoOWtwbCVTRZvXyED")
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
)

type Claims struct {
	UserId int64
	jwt2.RegisteredClaims
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (j *Handler) Refresh(tokenStr string) (string, error) {
	var claims Claims
	token, err := jwt2.ParseWithClaims(tokenStr, &claims, func(token *jwt2.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return "", ErrTokenInvalid
	}
	return j.GenerateToken(claims.UserId)
}

func (j *Handler) GenerateToken(uid int64) (string, error) {
	token := jwt2.NewWithClaims(jwt2.SigningMethodHS512, &Claims{
		UserId: uid,
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(time.Hour * 24 * 2)),
		},
	})
	tokenString, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *Handler) GenerateRefreshToken(uid int64) (string, error) {
	token := jwt2.NewWithClaims(jwt2.SigningMethodHS512, &Claims{
		UserId: uid,
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	})
	tokenString, err := token.SignedString(refreshSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *Handler) ParseToken(tokenString string) (Claims, error) {
	var claims Claims
	token, err := jwt2.ParseWithClaims(tokenString, &claims, func(token *jwt2.Token) (interface{}, error) {
		return tokenSecret, nil
	})
	if err != nil || !token.Valid {
		return Claims{}, ErrTokenInvalid
	}
	return claims, nil
}
