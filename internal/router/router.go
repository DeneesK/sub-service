package router

import (
	"time"

	_ "github.com/DeneesK/sub-service/api/docs"
	"github.com/DeneesK/sub-service/internal/model"
	"github.com/DeneesK/sub-service/internal/router/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type SubService interface {
	Create(sub *model.Subscription) error
	Get(id string) (*model.Subscription, error)
	List(userID string) ([]model.Subscription, error)
	Update(id string, upd *model.UpdateSubscription) error
	Delete(id string) error
	Aggregate(from, to time.Time, userID, service string) (int, error)
}

func NewRouter(timeOut time.Duration, subService SubService, log *zap.SugaredLogger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.NewLoggingMiddleware(log))
	r.Use(middleware.Timeout(timeOut * time.Second))

	h := NewSubscriptionHandler(subService, log)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/subs", h.Create)
		r.Get("/subs", h.List)
		r.Get("/subs/{id}", h.Get)
		r.Put("/subs/{id}", h.Update)
		r.Delete("/subs/{id}", h.Delete)
		r.Get("/subs/aggregate", h.Aggregate)
	})
	return r
}
