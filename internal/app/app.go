package app

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/DeneesK/sub-service/internal/model"
	"github.com/DeneesK/sub-service/internal/router"
	"go.uber.org/zap"
)

const shutdownTimeout = time.Second * 1

type SubService interface {
	Create(sub *model.Subscription) error
	Get(id string) (*model.Subscription, error)
	List(userID string) ([]model.Subscription, error)
	Update(id string, upd *model.Subscription) error
	Delete(id string) error
	Aggregate(from, to time.Time, userID, service string) (int, error)
}

type APP struct {
	srv        *http.Server
	log        *zap.SugaredLogger
	subService SubService
}

func NewApp(addr string, timeOut time.Duration, log *zap.SugaredLogger, subService SubService) *APP {
	r := router.NewRouter(timeOut, subService, log)
	s := http.Server{
		Addr:    addr,
		Handler: r,
	}
	return &APP{
		srv:        &s,
		log:        log,
		subService: subService,
	}
}

func (a *APP) Run() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	defer stop()

	a.log.Infoln("starting application, server listening on", a.srv.Addr)

	go func() {
		err := a.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.log.Fatalf("failed to start server: %s", err)
		}
	}()

	<-ctx.Done()

	a.log.Infoln("application shutdown process...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := a.srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Error during shutdown: %s", err)
	}
	<-shutdownCtx.Done()
	a.log.Infoln("application and server gracefully stopped")
}
