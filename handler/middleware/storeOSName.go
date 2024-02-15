package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type ctxKeyOS struct{}
var CtxKeyOS = ctxKeyOS{}

func StoreOSName(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		userAgent := r.UserAgent()
		ua := useragent.Parse(userAgent)

		ctx := context.WithValue(r.Context(), CtxKeyOS, ua.OS)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}