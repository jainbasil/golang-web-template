package rest

import (
	"github.com/labstack/echo/v4"
	"golang-web-template/internal/app/rest/handlers"
)

func (s *Server) initRoutes() {
	s.Logger.Info("Initializing Routes")
	initHealthCheckRoutes(s)
	initSampleApiRoutes(s)
}

func initSampleApiRoutes(s *Server, m ...echo.MiddlewareFunc) {
	sampleApiHandler := handlers.NewSampleApiHandler()
	s.POST("/samples", sampleApiHandler.DoSomething)
}

func initHealthCheckRoutes(s *Server, m ...echo.MiddlewareFunc) {
	healthCheckHandler := handlers.NewHealthCheckHandler()
	s.GET("/health", healthCheckHandler.Status)
}
