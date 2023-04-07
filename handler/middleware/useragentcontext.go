package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/mileusna/useragent"
)

type ctxKey string

const (
	userAgentKey ctxKey = "userAgent"
)

// UserAgentContext returns a middleware that adds user agent to context
func UserAgentContext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("user-agent: ", r.UserAgent())
		ua := useragent.Parse(r.UserAgent())
		ctx := r.Context()
		ctx = context.WithValue(ctx, userAgentKey, ua)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
