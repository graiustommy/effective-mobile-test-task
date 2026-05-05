package handler

import (
	"encoding/json"
	"net/http"

	"effective-task/internal/entity"
	apperrors "effective-task/internal/errors"
	"effective-task/internal/service"

	"go.uber.org/zap"
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

func NewHandler(logger *zap.Logger, service service.ServiceInterface) HandlerInterface {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type Handler struct {
	service service.ServiceInterface
	logger  *zap.Logger
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling create subscription request")
	var CreateRequest entity.Subscription
	err := json.NewDecoder(r.Body).Decode(&CreateRequest)
	if err != nil {
		h.logger.Warn("failed to parse request body", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	id, err := h.service.Create(r.Context(), &CreateRequest)
	if err != nil {
		h.logger.Error("failed to create subscription", zap.Error(err), zap.String("user_id", CreateRequest.UserID.String()), zap.String("service_name", CreateRequest.ServiceName))
		ResponseError(w, err)
		return
	}
	h.logger.Info("subscription created successfully", zap.String("subscription_id", id.String()), zap.String("user_id", CreateRequest.UserID.String()))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		h.logger.Error("failed to send response", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
}

func (h *Handler) ReadBySubscriptionID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling read subscription request")
	var SubscriptionID string
	err := json.NewDecoder(r.Body).Decode(&SubscriptionID)
	if err != nil {
		h.logger.Warn("failed to parse request body", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	Subscription, err := h.service.ReadBySubscriptionID(r.Context(), SubscriptionID)
	if err != nil {
		h.logger.Error("failed to read subscription", zap.Error(err), zap.String("subscription_id", SubscriptionID))
		ResponseError(w, err)
		return
	}
	h.logger.Info("subscription retrieved successfully", zap.String("subscription_id", SubscriptionID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(Subscription); err != nil {
		h.logger.Error("failed to send response", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling update subscription request")
	var CreateRequest entity.Subscription
	err := json.NewDecoder(r.Body).Decode(&CreateRequest)
	if err != nil {
		h.logger.Warn("failed to parse request body", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	err = h.service.Update(r.Context(), &CreateRequest)
	if err != nil {
		h.logger.Error("failed to update subscription", zap.Error(err), zap.String("subscription_id", CreateRequest.ID.String()))
		ResponseError(w, err)
		return
	}
	h.logger.Info("subscription updated successfully", zap.String("subscription_id", CreateRequest.ID.String()))
	status := http.StatusAccepted
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		h.logger.Error("failed to send response", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling delete subscription request")
	var subID string
	err := json.NewDecoder(r.Body).Decode(&subID)
	if err != nil {
		h.logger.Warn("failed to parse request body", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	err = h.service.Delete(r.Context(), subID)
	if err != nil {
		h.logger.Error("failed to delete subscription", zap.Error(err), zap.String("subscription_id", subID))
		ResponseError(w, err)
		return
	}
	h.logger.Info("subscription deleted successfully", zap.String("subscription_id", subID))
	status := http.StatusAccepted
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(status); err != nil {
		h.logger.Error("failed to send response", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
}

func (h *Handler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling list subscriptions request")
	var userID string
	err := json.NewDecoder(r.Body).Decode(&userID)
	if err != nil {
		h.logger.Warn("failed to parse request body", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	list, err := h.service.ListByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to list subscriptions", zap.Error(err), zap.String("user_id", userID))
		ResponseError(w, err)
		return
	}
	h.logger.Info("subscriptions listed successfully", zap.String("user_id", userID), zap.Int("count", len(list)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		h.logger.Error("failed to send response", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
	return
}

func (h *Handler) CountByUserID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling count by user id request")
	var request entity.CountPriceByUserID
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Warn("failed to parse request body", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	list, err := h.service.CountByUserID(r.Context(), &request)
	if err != nil {
		h.logger.Error("failed to count by user id", zap.Error(err), zap.String("user_id", request.UserID))
		ResponseError(w, err)
		return
	}
	h.logger.Info("count by user id retrieved successfully", zap.String("user_id", request.UserID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		h.logger.Error("failed to send response", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindInternal, "failed to send data", err))
	}
	return
}

func (h *Handler) CountByServiceName(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling count by service name request")
	var request entity.CountPriceByServiceName
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Warn("failed to parse request body", zap.Error(err))
		ResponseError(w, apperrors.Wrap(apperrors.ErrKindValidation, "failed to parse request body", err))
		return
	}
	defer r.Body.Close()
	list, err := h.service.CountByServiceName(r.Context(), &request)
	if err != nil {
		h.logger.Error("failed to count by service name", zap.Error(err), zap.String("service_name", request.ServiceName))
		ResponseError(w, err)
		return
	}
	h.logger.Info("count by service name retrieved successfully", zap.String("service_name", request.ServiceName))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		h.logger.Error("failed to send response", zap.Error(err))
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
