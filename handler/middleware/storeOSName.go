package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type ctxKeyOS struct{}

func StoreOSName(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userAgent := useragent.Parse(r.UserAgent())

		ctx := context.WithValue(r.Context(), ctxKeyOS{}, userAgent.OS)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}