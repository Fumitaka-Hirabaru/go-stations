package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	svc := service.NewTODOService(todoDB)
	mux.HandleFunc("/todos", handler.NewTODOHandler(svc).ServeHTTP)
	mux.HandleFunc("/healthz", handler.NewHealthzHandler().ServeHTTP)
	mux.HandleFunc("/panic", handler.NewPanicHandler().ServeHTTP)
	mux.HandleFunc("/do-panic", middleware.Recovery(handler.NewPanicHandler()).ServeHTTP)
	mux.HandleFunc("/user-agent", middleware.UserAgentContext(handler.NewUserAgentHandler()).ServeHTTP)
	mux.HandleFunc("/access-log", middleware.UserAgentContext(middleware.AccessLogMiddleware(handler.NewAccessLogHandler())).ServeHTTP)
	return mux
}
