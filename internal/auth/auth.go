package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func HandleWithAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			req, wr := setAuth(r, w)
			next.ServeHTTP(wr, req)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		var req *http.Request
		var wr http.ResponseWriter
		if err != nil || !token.Valid {
			req, wr = setAuth(r, w)
		} else {
			userID, _ := GetUserIDFromToken(cookie.Value)
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			req = r.WithContext(ctx)
		}

		next.ServeHTTP(wr, req)
	})
}

func setAuth(r *http.Request, w http.ResponseWriter) (*http.Request, http.ResponseWriter) {
	expTime := time.Now().Add(tokenExp)
	userID := uuid.New().String()

	tokenString, err := GenerateJWT(userID, expTime)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return nil, nil
	}

	ctx := context.WithValue(r.Context(), userIDKey, userID)
	r = r.WithContext(ctx)

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  expTime,
		Path:     "/",
		HttpOnly: true,
	})

	w.Header().Set("Authorization", "Bearer "+tokenString)

	return r, w
}
