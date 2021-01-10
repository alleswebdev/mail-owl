package main

import (
	"fmt"
	"github.com/alleswebdev/mail-owl/cmd/scheduler/app"
	"github.com/alleswebdev/mail-owl/internal/services/broker/rabbitmq"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := App.New()
	defer app.Storage.Db.Close()
	defer app.Logger.Sync()

	app.Logger.Infof("scheduler is running,: %#v", app.Config)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Kill)
		errs <- fmt.Errorf("%s", <-c)
	}()

	for i := 1; i < 6; i++ {
		err := app.Broker.Subscribe(rabbitmq.SchedulerQueue, app.SchedulerHandler)

		if err != nil {
			errs <- err
		}
	}

	app.Logger.Fatalw("getting error on the errors channel", "err", <-errs)
}
