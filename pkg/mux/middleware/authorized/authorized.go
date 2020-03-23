package authorized

import (
	"context"
	"log"
	"net/http"

	"github.com/banch0/crud2/pkg/jwt"
)

// Authorized ...
func Authorized(roles []string, payload func(ctx context.Context) interface{}) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			auth := payload(request.Context()).(*jwt.JWTPayload)
			for _, role := range roles {
				for _, r := range auth.Roles {
					if role == r {
						log.Printf("access granted %v %v", roles, auth)
						next(writer, request)
						return
					}
				}
			}
			http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	}
}
