package middleware

import (
	"net/http"
)

func Log(h http.Handler) http.Handler {
fn := func(w http.ResponseWriter, r *http.Request) {
	// 処理書くとこ

	h.ServeHTTP(w, r)
}
return http.HandlerFunc(fn)
}