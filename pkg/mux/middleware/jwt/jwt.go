package jwt

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	jwtcore "github.com/banch0/crud2/pkg/jwt"
)

type contextKey string

var payloadContextKey = contextKey("jwt")

// JWT ...
func JWT(payloadType reflect.Type, secret jwtcore.Secret) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {

			header := request.Header.Get("Authorization")
			if header == "" {
				http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			var token string

			if len(header) > len("Bearer ") {
				token = header[len("Bearer "):]
			}

			ok, err := jwtcore.Verify(token, secret)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				log.Println("JWT Bad Request: ", err)
				return
			}

			if !ok {
				log.Println("401 Unautorized")
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			payload := reflect.New(payloadType).Interface()

			_, err = jwtcore.Decode(payload, token)
			if err != nil {
				log.Println("JWT Bad Request: ", err)
				http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			ok, err = jwtcore.IsNotExpired(payload, time.Now())
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				log.Println("JWT Bad Request: ", err)
				return
			}

			if !ok {
				log.Println("Unautorized Token Expired")
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Add to context payload with key "jwt"
			ctx := context.WithValue(request.Context(), payloadContextKey, payload)
			next(writer, request.WithContext(ctx))
		}
	}
}

// FromContext ...
func FromContext(ctx context.Context) (payload interface{}) {
	payload = ctx.Value(payloadContextKey)
	return
}

// IsContextNonEmpty ...
func IsContextNonEmpty(ctx context.Context) bool {
	return nil != ctx
}
