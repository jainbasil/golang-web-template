package main

import (
	"golang-web-template/internal"
	"golang-web-template/internal/app"
	"golang-web-template/internal/app/events/handlers"
	"golang-web-template/internal/app/rest"
	"golang-web-template/internal/config"
	"golang-web-template/internal/mq/rabbitmq"
)

func main() {
	cfg := config.LoadConfig()
	appContext := internal.InitAppContext(cfg)
	apiService := rest.NewServer(cfg, appContext)

	// initialize the event handler for messages in queues
	// Note: You can register any number of event handlers, and it will
	// handle based on the key specified in the handler.
	sampleEventHandler := handlers.NewSampleEventHandler(appContext)

	// Uncomment the following code if you need rabbitmq consumer and add it to Run()
	rmqConsumer := rabbitmq.NewConsumer(cfg, appContext)
	rmqConsumer.RegisterCloudEventHandler(sampleEventHandler)

	// Uncomment the following code if you need pubsub consumer and add it to Run()
	//pubsubConsumer := pubsub.NewConsumer(cfg, appContext)
	//pubsubConsumer.RegisterCloudEventHandler(sampleEventHandler)

	// add runnable services to Run method for starting
	app.Run(appContext, apiService)
}
