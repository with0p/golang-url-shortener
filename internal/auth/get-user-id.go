package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func GetUserIDFromToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}

func GetUserIDFromCtx(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", errors.New("no user id")
	}

	return userID, nil
}

func GetUserIDFromCookie(r *http.Request) (string, error) {
	cookie, cookieErr := r.Cookie("auth_token")
	if cookieErr != nil {
		return "", cookieErr
	}

	userID, err := GetUserIDFromToken(cookie.Value)

	if err != nil {
		return "", err
	}

	return userID, nil
}
