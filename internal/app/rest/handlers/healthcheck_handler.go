package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthCheckHandler struct {
	*BaseHandler
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) Status(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
}
