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
)

type Producer struct {
	host              string
	port              string
	user              string
	password          string
	exchange          string
	queue             string
	appContext        *internal.AppContext
	connection        *amqp.Connection
	channel           *amqp.Channel
	shutdownRequested bool
}

func NewProducer(cfg *config.AppConfig, appContext *internal.AppContext) *Producer {
	producer := &Producer{
		host:       cfg.RabbitMQ.Host,
		port:       cfg.RabbitMQ.Port,
		user:       cfg.RabbitMQ.Username,
		password:   cfg.RabbitMQ.Password,
		exchange:   cfg.RabbitMQ.Exchange,
		queue:      cfg.RabbitMQ.Queue,
		appContext: appContext,
	}
	err := producer.Connect()
	if err != nil {
		appContext.Logger.Sugar().Panic("(producer):: error connecting to rabbitmq", err)
	}
	return producer
}

func (p *Producer) PublishMessage(event domain.CloudEvent, attributes map[string]string) error {
	message, err := json.Marshal(&event)
	if err != nil {
		return err
	}
	err = p.channel.PublishWithContext(
		context.Background(),
		"",
		p.queue,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Connect() error {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", p.user, p.password, p.host, p.port)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	p.connection = conn
	p.channel = ch

	//TODO enhance by declaring the exchange, queue etc.

	go p.monitorConnection()
	return nil
}

func (p *Producer) Run() {
	// implemented for interface compliance and intentionally
	// left empty as there is nothing to be run
}

func (p *Producer) Stop(ctx context.Context) {
	p.appContext.Logger.Info("(producer):: shutting down")
	p.shutdownRequested = true
	if err := p.channel.Close(); err != nil {
		p.appContext.Logger.Warn("(producer):: error closing rabbitmq channel")
		return
	}

	if err := p.connection.Close(); err != nil {
		p.appContext.Logger.Warn("(producer):: error closing rabbitmq connection")
		return
	}
}

func (p *Producer) monitorConnection() {
	p.appContext.Logger.Sugar().Debug("(producer):: monitoring rabbitmq connection")
	closeCh := make(chan *amqp.Error)
	p.connection.NotifyClose(closeCh)

	select {
	case <-closeCh:
		if p.shutdownRequested {
			return
		}
		for {
			p.appContext.Logger.Sugar().Debug("(producer):: attempting reconnection to rabbitmq")
			err := p.Connect()
			if err == nil {
				return
			}
			p.appContext.Logger.Sugar().Error("(producer):: reconnecting to rabbitmq failed, retrying after 5 seconds!", zap.Error(err))
			time.Sleep(5 * time.Second)
		}
	}
}
