package rest

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang-web-template/internal"
	"golang-web-template/internal/config"
	"net/http"
)

// Server exposes an API service based on echo.Echo. All the required dependencies
// like database connections, services etc. are injected into it through the internal.AppContext.
type Server struct {
	*echo.Echo
	port       string
	appContext *internal.AppContext
}

// NewServer initializes and return an instance of rest.Server using the application configuration
// (config.AppConfig) and context (internal.AppContext). It also initializes the middlewares and routes
// exposed by the Server.
func NewServer(cfg *config.AppConfig, appContext *internal.AppContext) *Server {
	s := &Server{
		Echo:       echo.New(),
		port:       cfg.Port,
		appContext: appContext,
	}

	s.initMiddlewares()
	s.initRoutes()
	return s
}

// Run method starts the server at the configured port. This is an implementation of Run method signature defined
// in app.Runnable interface.
func (s *Server) Run() {
	port := fmt.Sprintf(":%s", s.port)
	s.appContext.Logger.Sugar().Infof("starting zoko-payments api service at address %s ", port)
	if err := s.Start(port); err != nil && err != http.ErrServerClosed {
		s.appContext.Logger.Sugar().Panic("error occurred, shutting down server", err)
	}
}

// Stop method gracefully stops the rest.Server by invoking Shutdown method of Echo.
// This is an implementation of Stop method signature defined in app.Runnable interface.
func (s *Server) Stop(ctx context.Context) {
	s.appContext.Logger.Info("shutting down http server")
	if err := s.Shutdown(ctx); err != nil {
		s.appContext.Logger.Sugar().Fatal(err)
	}
}
