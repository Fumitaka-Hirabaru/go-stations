package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func applyMiddleware(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	svc := service.NewTODOService(todoDB)
	middlewares := []func(http.Handler) http.Handler{
		middleware.UserAgentContext,
		middleware.AccessLogMiddleware,
	}
	mux.HandleFunc("/todos", handler.NewTODOHandler(svc).ServeHTTP)
	mux.HandleFunc("/healthz", handler.NewHealthzHandler().ServeHTTP)
	mux.HandleFunc("/panic", handler.NewPanicHandler().ServeHTTP)
	mux.HandleFunc("/do-panic", middleware.Recovery(handler.NewPanicHandler()).ServeHTTP)
	mux.HandleFunc("/user-agent", middleware.UserAgentContext(handler.NewUserAgentHandler()).ServeHTTP)
	mux.HandleFunc("/access-log", applyMiddleware(handler.NewAccessLogHandler(), middlewares...).ServeHTTP)
	return mux
}
