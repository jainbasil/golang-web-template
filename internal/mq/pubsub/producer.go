package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"zoko-go-web-template/internal"
	"zoko-go-web-template/internal/config"
	"zoko-go-web-template/internal/domain"
)

type Producer struct {
	client *pubsub.Client
	topic  *pubsub.Topic
	logger *zap.Logger
}

func NewProducer(cfg *config.AppConfig, appContext *internal.AppContext) *Consumer {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, cfg.GcpConfig.ProjectID)
	if err != nil {
		// TODO return error as well for better handling
		return nil
	}

	topic := client.Topic(cfg.GcpConfig.PubsubConfig.Topic)

	return &Consumer{
		client: client,
		topic:  topic,
		logger: appContext.Logger,
	}
}

func (p *Producer) PublishMessage(event domain.CloudEvent, attributes map[string]string) error {
	ctx := context.Background()
	payload, _ := json.Marshal(event)
	result := p.topic.Publish(ctx, &pubsub.Message{
		Data:       payload,
		Attributes: attributes,
	})
	_, err := result.Get(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Connect() error {
	return nil
}

func (p *Producer) Run() {
	// implemented for interface compliance and intentionally
	// left empty as there is nothing to be run
}

func (p *Producer) Stop(ctx context.Context) {
	err := p.client.Close()
	if err != nil {
		p.logger.Warn("(producer):: error closing pubsub client")
	}
}
