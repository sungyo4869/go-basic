package middleware

import (
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
fn := func(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {	
			http.Error(w, "Internal Server Error", 500)
		}
	}()

	h.ServeHTTP(w, r)
}
return http.HandlerFunc(fn)
}