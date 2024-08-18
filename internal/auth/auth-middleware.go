package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			setAuth(r, w)
			next.ServeHTTP(w, r)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})

		if err != nil || !token.Valid {
			setAuth(r, w)
		} else {
			userId, _ := GetUserIdFromToken(cookie.Value)
			ctx := context.WithValue(r.Context(), "userId", userId)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func setAuth(r *http.Request, w http.ResponseWriter) {
	expTime := time.Now().Add(TOKEN_EXP)
	userId := uuid.New().String()

	tokenString, err := GenerateJWT(userId, expTime)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(r.Context(), "userId", userId)
	r = r.WithContext(ctx)

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  expTime,
		Path:     "/",
		HttpOnly: true,
	})

	w.Header().Set("Authorization", "Bearer "+tokenString)
}
