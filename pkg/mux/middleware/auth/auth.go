package auth

import (
	"context"
	"log"
	"net/http"
)

type contextKey string

var tokenContextKey = contextKey("jwt")

// Auth ...
func Auth() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			log.Println("Auth service")
			token := req.Header.Get("Authorization")
			if token == "" {
				next(res, req)
				return
			}

			next(res, req)

			ctx := context.WithValue(
				req.Context(),
				tokenContextKey,
				token,
			)
			next(res, req.WithContext(ctx))
		}
	}
}

// FromContext ...
func FromContext(ctx context.Context) (token string, ok bool) {
	token, ok = ctx.Value(tokenContextKey).(string)
	return
}
