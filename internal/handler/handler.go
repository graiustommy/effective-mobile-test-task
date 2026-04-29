package handler

import (
	"net/http"

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
	router chi.Router
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	return
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
