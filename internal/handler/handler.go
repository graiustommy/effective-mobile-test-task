package handler

import (
	"encoding/json"
	"net/http"

	"effective-task/internal/entity"
	apperrors "effective-task/internal/errors"
	"effective-task/internal/service"
)

type HandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	ReadBySubscriptionID(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	ListByUserID(w http.ResponseWriter, r *http.Request)
	CountByUserID(w http.ResponseWriter, r *http.Request)
	CountByServiceName(w http.ResponseWriter, r *http.Request)
}

func NewHandler(service service.ServiceInterface) HandlerInterface {
	return &Handler{
		service: service,
	}
}

type Handler struct {
	service service.ServiceInterface
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var CreateRequest entity.Subscription
	err := json.NewDecoder(r.Body).Decode(&CreateRequest)
	if err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	id, err := h.service.Create(r.Context(), &CreateRequest)
	if err != nil {
		ResponseError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
}

func (h *Handler) ReadBySubscriptionID(w http.ResponseWriter, r *http.Request) {
	var SubscriptionID string
	err := json.NewDecoder(r.Body).Decode(&SubscriptionID)
	if err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	Subscription, err := h.service.ReadBySubscriptionID(r.Context(), SubscriptionID)
	if err != nil {
		ResponseError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(Subscription); err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var CreateRequest entity.Subscription
	err := json.NewDecoder(r.Body).Decode(&CreateRequest)
	if err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	err = h.service.Update(r.Context(), &CreateRequest)
	if err != nil {
		ResponseError(w, err)
		return
	}
	status := http.StatusAccepted
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var subID string
	err := json.NewDecoder(r.Body).Decode(&subID)
	if err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	err = h.service.Delete(r.Context(), subID)
	if err != nil {
		ResponseError(w, err)
		return
	}
	status := http.StatusAccepted
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(status); err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
}

func (h *Handler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	var userID string
	err := json.NewDecoder(r.Body).Decode(&userID)
	if err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	list, err := h.service.ListByUserID(r.Context(), userID)
	if err != nil {
		ResponseError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
	return
}

func (h *Handler) CountByUserID(w http.ResponseWriter, r *http.Request) {
	var request entity.CountPriceByUserID
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	list, err := h.service.CountByUserID(r.Context(), &request)
	if err != nil {
		ResponseError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
	return
}

func (h *Handler) CountByServiceName(w http.ResponseWriter, r *http.Request) {
	var request entity.CountPriceByServiceName
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	list, err := h.service.CountByServiceName(r.Context(), &request)
	if err != nil {
		ResponseError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
	return
}

func ResponseError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	statusCode := apperrors.StatusCodeForError(err)
	w.WriteHeader(statusCode)

	errResponse := apperrors.ToErrorResponse(err)
	json.NewEncoder(w).Encode(errResponse)
}
