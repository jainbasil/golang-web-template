package handlers

import (
	"go.uber.org/zap"
	"golang-web-template/internal"
	"golang-web-template/internal/domain"
)

// SampleEventHandler implements mq.CloudEventHandler interface and process the cloud events
// consumed by the application from rabbitmq or gcp pub-sub
type SampleEventHandler struct {
	logger *zap.Logger
}

func NewSampleEventHandler(ctx *internal.AppContext) *SampleEventHandler {
	return &SampleEventHandler{
		logger: ctx.Logger,
	}
}

func (h *SampleEventHandler) ProcessEvent(e domain.CloudEvent) error {
	h.logger.Sugar().Info("processing event", e.Data)
	return nil
}

func (h *SampleEventHandler) GetCloudEventSourceKey() string {
	return "com.github"
}
