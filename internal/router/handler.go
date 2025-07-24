package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DeneesK/sub-service/internal/model"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	svc SubService
	log *zap.SugaredLogger
}

func NewSubscriptionHandler(svc SubService, log *zap.SugaredLogger) *SubscriptionHandler {
	return &SubscriptionHandler{
		svc: svc,
		log: log,
	}
}

// CreateSubscription
// @Summary Create subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.Subscription true
// @Success 201 {object} model.Subscription
// @Router /api/v1/subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.Subscription
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.Create(&req); err != nil {
		h.log.Error("create error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// CetSubscription
// @Summary Get subscription
// @Description Get subscription by its id
// @Tags subscriptions
// @Produce json
// @Param id query string true
// @Param subscription body model.Subscription true
// @Success 200 {object} model.Subscription
// @Router /api/v1/subs/{id} [get]
func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "must provide id", http.StatusBadRequest)
		return
	}
	sub, err := h.svc.Get(id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

// ListSubscription
// @Summary Get list of subscriptions
// @Description Get list of subscriptions by user_id if user_id is empty returns all subs
// @Tags subscriptions
// @Produce json
// @Param user_id query string false
// @Success 200 {array} model.Subscription
// @Router /api/v1/subs [get]
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	subs, err := h.svc.List(userID)
	if err != nil {
		h.log.Error("list error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subs)
}

// UpdateSubscription
// @Summary Update subscription
// @Description Update subscription by its id
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body model.Subscription true
// @Success 200 {object} model.Subscription
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /api/v1/subs/{id} [put]
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	var req model.Subscription
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.ID = id
	if err := h.svc.Update(id, &req); err != nil {
		h.log.Error("update error", zap.Error(err))
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

// DeleteSubscription
// @Summary Delete subscription
// @Description Delete subscription by its id
// @Tags subscriptions
// @Produce plain
// @Param id path string true "Subscription ID"
// @Success 204 {string} string "No Content"
// @Failure 404 {string} string "Not Found"
// @Router /api/v1/subs/{id} [delete]
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	if err := h.svc.Delete(id); err != nil {
		h.log.Error("delete error", zap.Error(err))
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AggregateSubscription
// @Summary Aggregate subscriptions cost
// @Description Sum prices between dates, optional filters user_id & service_name
// @Tags subscriptions
// @Produce json
// @Param from query string true "Start month-year" example(2025-01)
// @Param to query string true "End month-year"   example(2025-07)
// @Param user_id query string false
// @Param service_name query string false
// @Success 200 {object} map[string]int
// @Router /api/v1/subs/aggregate [get]
func (h *SubscriptionHandler) Aggregate(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")

	from, err := time.Parse("2006-01", fromStr)
	if err != nil {
		http.Error(w, "invalid from", http.StatusBadRequest)
		return
	}
	to, err := time.Parse("2006-01", toStr)
	if err != nil {
		http.Error(w, "invalid to", http.StatusBadRequest)
		return
	}

	sum, err := h.svc.Aggregate(from, to.AddDate(0, 1, -1), userID, serviceName)
	if err != nil {
		h.log.Error("aggregate error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"total": sum})
}
