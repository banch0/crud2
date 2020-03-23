package auth

import (
	"context"
	"net/http"
)

type contextKey string

var tokenContextKey = contextKey("jwt")

// Auth ...
func Auth(ctx func(ctx context.Context) bool) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			if !ctx(request.Context()) {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			next(writer, request)
		}
	}
}

// FromContext ...
func FromContext(ctx context.Context) (token string, ok bool) {
	token, ok = ctx.Value(tokenContextKey).(string)
	return
}
