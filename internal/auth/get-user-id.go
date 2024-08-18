package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func GetUserIdFromToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}

func GetUserIdFromCtx(ctx context.Context) (string, error) {
	fmt.Println(ctx.Value("userId"))
	userID, ok := ctx.Value("userId").(string)
	if !ok {
		return "", errors.New("no user id")
	}

	return userID, nil
}
