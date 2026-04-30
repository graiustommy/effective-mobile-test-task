package handler

import (
	"encoding/json"
	"net/http"

	"effective-task/internal/entity"
	"effective-task/internal/service"

	chi "github.com/go-chi/chi/v5"
)

type HandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	router  *chi.Router
	service service.ServiceInterface
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var CreateRequest entity.Subscription
	err := json.NewDecoder(r.Body).Decode(&CreateRequest)
	if err != nil {
		ResponseError(w, "Invalid JSON format")
		return
	}
	defer r.Body.Close()
	err = h.service.Create(r.Context(), &CreateRequest)
	if err != nil {
		ResponseError(w, err.Error())
		return
	}
}

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	return
}

func ResponseError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(entity.ErrorResponse{
		Error: msg,
	})
}
