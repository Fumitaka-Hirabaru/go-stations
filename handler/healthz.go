package handler

import (
	"log"
	"encoding/json"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{
	Endpoint string
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{Endpoint: "/healthz"}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	healthzResponse := &model.HealthzResponse{Message: "OK"}
	err := json.NewEncoder(w).Encode(healthzResponse)
	if err != nil {
		log.Println(err)
	}
}
