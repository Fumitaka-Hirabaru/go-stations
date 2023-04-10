package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AccessLog struct {
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency"`
	Path      string    `json:"path"`
	OS        string    `json:"os"`
}

func AccessLogMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		end := time.Now()
		latency := end.Sub(start).Milliseconds()

		path := r.URL.Path
		ua := r.Context().Value(userAgentKey{}).(string)
		al := AccessLog{
			Timestamp: start,
			Latency:   latency,
			Path:      path,
			OS:        ua,
		}
		b, err := json.Marshal(al)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(b))
	}
	return http.HandlerFunc(fn)
}
