package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DeneesK/sub-service/internal/model"
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
// @Router /api/v1/subscriptions/aggregate [get]
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
