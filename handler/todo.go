package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "POST":
		createTODORequest := &model.CreateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(createTODORequest); err != nil {
			log.Println(err)
			return
		}

		if createTODORequest.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		createTODOResponse, err := h.Create(ctx, createTODORequest)
		if err != nil {
			log.Println(err)
			return
		}

		json.NewEncoder(w).Encode(createTODOResponse)

	case "PUT":
		updateTODORequest := model.UpdateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(&updateTODORequest); err != nil {
			log.Println(err)
			return
		}

		if updateTODORequest.ID == 0 || updateTODORequest.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updateTODOResponse, err := h.Update(r.Context(), &updateTODORequest)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(updateTODOResponse)

	case "GET":
		query := r.URL.Query()
		prev_id, _ := strconv.ParseInt(query.Get("prev_id"), 10, 64)
		size, _ := strconv.ParseInt(query.Get("size"), 10, 64)
		readTODORequest := &model.ReadTODORequest{
			PrevID: prev_id,
			Size:   size,
		}

		ctx := r.Context()
		readTODOResponse, err := h.Read(ctx, readTODORequest)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(readTODOResponse)

	case "DELETE":
		deleteTODORequest := &model.DeleteTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(deleteTODORequest); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(deleteTODORequest.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		deleteTODOResponse, err := h.Delete(ctx, deleteTODORequest)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(deleteTODOResponse)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	return &model.CreateTODOResponse{TODO: *todo}, err
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	res := &model.ReadTODOResponse{TODOs: make([]model.TODO, 0)}
	for _, todo := range todos {
		res.TODOs = append(res.TODOs, *todo)
	}
	return res, err
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	return &model.UpdateTODOResponse{TODO: *todo}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}
