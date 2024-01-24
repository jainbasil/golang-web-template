package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-web-template/internal"
)

// Runnable interface needs to be implemented by all services that has to be started as part of the main function.
// This interface is used to run/stop the implementing services as an independent concurrent thread of control, or goroutine,
// within the same address space. See Run() for implementation details.
type Runnable interface {
	Run()
	Stop(ctx context.Context)
}

func Run(appContext *internal.AppContext, services ...Runnable) {
	for _, service := range services {
		go service.Run()
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	var e = <-exit
	appContext.Logger.Sugar().Infof("%s signal received, quitting.. ", e)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, service := range services {
		service.Stop(ctx)
	}

	// uncomment the following if you have a db repository in your app context that needs to be disconnected from.
	//if err := appContext.DbRepository.Disconnect(); err != nil {
	//	appContext.Logger.Sugar().Error(err)
	//}
	_ = appContext.Logger.Sync()
}
