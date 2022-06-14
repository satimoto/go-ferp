package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/edjumacator/chi-prometheus"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/satimoto/go-datastore/pkg/util"
)

type Rest interface {
	Handler() *chi.Mux
	StartRest(context.Context, *sync.WaitGroup)
}

type RestService struct {
	*http.Server
}

func NewRest() Rest {
	return &RestService{}
}

func (rs *RestService) Handler() *chi.Mux {
	router := chi.NewRouter()

	// Set middleware
	router.Use(render.SetContentType(render.ContentTypeJSON), middleware.RedirectSlashes, middleware.Recoverer)
	router.Use(middleware.Timeout(120 * time.Second))
	router.Use(chiprometheus.NewMiddleware("ferp"))

	router.Mount("/health", rs.mountHealth())
	router.Mount("/metrics", promhttp.Handler())

	return router
}

func (rs *RestService) StartRest(ctx context.Context, waitGroup *sync.WaitGroup) {
	if rs.Server == nil {
		rs.Server = &http.Server{
			Addr:    fmt.Sprintf(":%s", os.Getenv("REST_PORT")),
			Handler: rs.Handler(),
		}
	}

	log.Printf("Starting Rest service")
	waitGroup.Add(1)

	go rs.listenAndServe()

	go func() {
		<-ctx.Done()
		log.Printf("Shutting down Rest service")

		rs.shutdown()

		log.Printf("Rest service shut down")
		waitGroup.Done()
	}()
}

func (rs *RestService) listenAndServe() {
	err := rs.Server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		util.LogOnError("FERP001", "Error in Rest service", err)
	}
}

func (rs *RestService) shutdown() {
	timeout := util.GetEnvInt32("SHUTDOWN_TIMEOUT", 20)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	err := rs.Server.Shutdown(ctx)

	if err != nil {
		util.LogOnError("FERP002", "Error shutting down Rest service", err)
	}
}
