package handlers

import "github.com/labstack/echo/v4"

type SampleApiHandler struct {
	*BaseHandler
}

func NewSampleApiHandler() *SampleApiHandler {
	return &SampleApiHandler{}
}

func (h *SampleApiHandler) DoSomething(c echo.Context) error {
	return nil
}
