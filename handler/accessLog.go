package handler

import "net/http"

type AccessLogHandler struct{}

func NewAccessLogHandler() *AccessLogHandler {
	return &AccessLogHandler{}
}

func (h *AccessLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
