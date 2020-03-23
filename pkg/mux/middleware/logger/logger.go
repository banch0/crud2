package logger

import (
	"log"
	"net/http"
)

// Logger ...
func Logger(prefix string) func(next http.HandlerFunc,
) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			log.Printf(
				"%s Method: %s, path: %s",
				prefix,
				req.Method,
				req.URL.Path,
			)
			next(res, req)
		}
	}
}
