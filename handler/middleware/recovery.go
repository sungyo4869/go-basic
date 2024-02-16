package middleware

import (
	"log"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Print(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}()

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
