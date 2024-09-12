package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/diianpro/tingerDog/transport/handler"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	// Docs contains documentation for swagger.
	_ "github.com/diianpro/tingerDog/docs"
)

type Server struct {
	httpServer   *http.Server
	publicRouter *chi.Mux

	handler *handler.Handler
}

func New(handler *handler.Handler) *Server {
	s := &Server{
		handler: handler,
	}

	router := chi.NewRouter()
	s.httpServer = &http.Server{
		IdleTimeout:       30 * time.Second,
		Handler:           router,
		ReadHeaderTimeout: 1 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	router.Get("/docs/*", httpSwagger.WrapHandler)
	router.Route("/", s.router)

	return s
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return s.httpServer.Close()
	}
	return nil
}

func (s *Server) Serve() {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return
	}
}

func (s *Server) router(r chi.Router) {
	r.Route("/list", func(r chi.Router) {
		r.Get("/user", s.handler.GetAllUsers)
	})
}
