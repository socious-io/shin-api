package auth

import (
	"errors"
	"shin/src/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID      string `json:"id"`
	Refresh bool   `json:"refresh"`
	jwt.RegisteredClaims
}

type SSOClaims struct {
	ID        *string `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(id string, refresh bool) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		ID:      id,
		Refresh: refresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.Secret))
}

func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Secret), nil
	})
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, errors.New("invalid token")
	} else if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("unknown claims type, cannot proceed")
}

func VerifySSOToken(tokenString string) (*SSOClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SSOClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.SSO.Secret), nil
	})
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, errors.New("invalid token")
	} else if claims, ok := token.Claims.(*SSOClaims); ok {
		return claims, nil
	}

	return nil, errors.New("unknown claims type, cannot proceed")
}

func GenerateFullTokens(id string) (map[string]any, error) {
	accessToken, err := GenerateToken(id, false)
	if err != nil {
		return nil, err
	}
	refreshToken, err := GenerateToken(id, true)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	}, nil
}
