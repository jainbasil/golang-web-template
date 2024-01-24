package handlers

import "go.uber.org/zap"

type BaseHandler struct {
	logger *zap.Logger
}

func NewBaseHandler(l *zap.Logger) *BaseHandler {
	return &BaseHandler{
		logger: l,
	}
}
