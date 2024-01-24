package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"

	"golang-web-template/internal"
	"golang-web-template/internal/config"
	"golang-web-template/internal/domain"
	"golang-web-template/internal/mq"
)

type Consumer struct {
	appContext         *internal.AppContext
	client             *pubsub.Client
	subscription       *pubsub.Subscription
	cloudEventHandlers map[string]mq.CloudEventHandler
}

func NewConsumer(cfg *config.AppConfig, appContext *internal.AppContext) *Consumer {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, cfg.GcpConfig.ProjectID)
	if err != nil {
		// TODO return error as well for better handling
		return nil
	}

	sub := client.Subscription(cfg.GcpConfig.PubsubConfig.Subscription)
	return &Consumer{
		client:       client,
		subscription: sub,
		appContext:   appContext,
	}
}

func (c *Consumer) Run() {
	if err := c.ConsumeMessages(); err != nil {
		// handle runtime error
	}
}

func (c *Consumer) Stop(ctx context.Context) {
	err := c.client.Close()
	if err != nil {
		c.appContext.Logger.Warn("(consumer):: error closing pubsub client")
	}
}

func (c *Consumer) ConsumeMessages() error {
	err := c.subscription.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		var cloudEvent domain.CloudEvent
		unmarshalErr := json.Unmarshal(m.Data, &cloudEvent)
		if unmarshalErr != nil {
			c.appContext.Logger.Sugar().Error("error unmarshalling message", unmarshalErr)
			m.Nack()
			return
		}
		if handler, ok := c.cloudEventHandlers[cloudEvent.Source]; ok {
			if err := handler.ProcessEvent(cloudEvent); err != nil {
				c.appContext.Logger.Sugar().Error("error processing cloud event", err)
				m.Nack()
				return
			}
		} else {
			c.appContext.Logger.Sugar().Error("unsupported message type")
		}
		m.Ack()
	})
	if !errors.Is(err, context.Canceled) {
		c.appContext.Logger.Sugar().Error("(consumer):: failed to receive from pubsub")
		return err
	}
	return nil
}

func (c *Consumer) RegisterCloudEventHandler(handler mq.CloudEventHandler) {
	key := handler.GetCloudEventSourceKey()
	c.appContext.Logger.Sugar().Debugf("registering handler for source: %s", key)
	c.cloudEventHandlers[key] = handler
}
