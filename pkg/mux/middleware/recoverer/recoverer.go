package recoverer

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recoverer ...
func Recoverer() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("%s", debug.Stack())
					http.Error(
						res,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError,
					)
				}
			}()
			next(res, req)
		}
	}
}
