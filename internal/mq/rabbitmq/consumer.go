package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"

	"golang-web-template/internal"
	"golang-web-template/internal/config"
	"golang-web-template/internal/domain"
	"golang-web-template/internal/mq"
)

type HandlerFunc func(e domain.CloudEvent) error

type Consumer struct {
	host               string
	port               string
	user               string
	password           string
	exchange           string
	queue              string
	appContext         *internal.AppContext
	connection         *amqp.Connection
	channel            *amqp.Channel
	cloudEventHandlers map[string]mq.CloudEventHandler
	shutdownRequested  bool
}

func NewConsumer(cfg *config.AppConfig, appContext *internal.AppContext) *Consumer {
	consumer := &Consumer{
		host:               cfg.RabbitMQ.Host,
		port:               cfg.RabbitMQ.Port,
		user:               cfg.RabbitMQ.Username,
		password:           cfg.RabbitMQ.Password,
		exchange:           cfg.RabbitMQ.Exchange,
		queue:              cfg.RabbitMQ.Queue,
		appContext:         appContext,
		cloudEventHandlers: make(map[string]mq.CloudEventHandler),
	}

	err := consumer.Connect()
	if err != nil {
		appContext.Logger.Sugar().Panic("error connecting to rabbitmq", err)
	}

	return consumer
}

func (c *Consumer) Run() {
	c.appContext.Logger.Info("(consumer):: starting up")
	forever := make(chan bool)

	if err := c.ConsumeMessages(); err != nil {
		c.appContext.Logger.Error("error consuming from rabbitmq", zap.String("err", err.Error()))
		//c.appContext.Logger.Error("exiting by emitting signal")
		//_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}

	c.appContext.Logger.Info("(consumer):: waiting for messages")
	<-forever
}

func (c *Consumer) Stop(ctx context.Context) {
	c.shutdownRequested = true
	c.appContext.Logger.Info("(consumer):: shutting down")
	if err := c.channel.Cancel(c.queue, false); err != nil {
		return
	}

	if err := c.connection.Close(); err != nil {
		return
	}
}

func (c *Consumer) ConsumeMessages() error {
	deliveries, err := c.channel.Consume(c.queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go c.processDeliveries(deliveries)
	return nil
}

func (c *Consumer) RegisterCloudEventHandler(handler mq.CloudEventHandler) {
	key := handler.GetCloudEventSourceKey()
	c.appContext.Logger.Sugar().Debugf("registering handler for source: %s", key)
	c.cloudEventHandlers[key] = handler
}

func (c *Consumer) Connect() error {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", c.user, c.password, c.host, c.port)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	err = ch.Qos(10, 0, false)
	if err != nil {
		return err
	}

	c.connection = conn
	c.channel = ch
	go c.monitorConnection()
	return nil
}

func (c *Consumer) processDeliveries(messages <-chan amqp.Delivery) {
	var cloudEvent domain.CloudEvent
	for msg := range messages {
		err := json.Unmarshal(msg.Body, &cloudEvent)
		if err != nil {
			c.appContext.Logger.Sugar().Error("error unmarshalling message", err)
		}
		if handler, ok := c.cloudEventHandlers[cloudEvent.Source]; ok {
			err := handler.ProcessEvent(cloudEvent)
			if err != nil {
				c.appContext.Logger.Sugar().Error("error processing cloud event", err)
			}
		} else {
			c.appContext.Logger.Sugar().Error("unsupported message type")
		}
	}
}

func (c *Consumer) monitorConnection() {
	c.appContext.Logger.Sugar().Debug("(consumer):: monitoring rabbitmq connection")
	closeCh := make(chan *amqp.Error)
	c.connection.NotifyClose(closeCh)

	select {
	case <-closeCh:
		if c.shutdownRequested {
			return
		}
		// retry connection only if shutdown was never requested
		for {
			c.appContext.Logger.Sugar().Debug("(consumer):: attempting reconnection to rabbitmq")
			err := c.Connect()
			if err == nil {
				return
			}
			c.appContext.Logger.Sugar().Error("(consumer):: reconnecting to rabbitmq failed, retrying after 5 seconds!", zap.Error(err))
			time.Sleep(5 * time.Second)
		}
	}
}
