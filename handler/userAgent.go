package handler

import (
	"net/http"
)

type UserAgentHandler struct{}

func NewUserAgentHandler() *UserAgentHandler {
	return &UserAgentHandler{}
}

func (h *UserAgentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
