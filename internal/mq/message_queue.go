package mq

import "golang-web-template/internal/domain"

type CloudEventHandler interface {
	GetCloudEventSourceKey() string
	ProcessEvent(event domain.CloudEvent) error
}

type MessageQueueConsumer interface {
	ConsumeMessages() error
	RegisterCloudEventHandler(handler CloudEventHandler)
}

type MessageQueueProducer interface {
	PublishMessage(event domain.CloudEvent, attributes map[string]string) error
}
