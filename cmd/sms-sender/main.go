package main

import (
	"fmt"
	"github.com/alleswebdev/mail-owl/cmd/sms-sender/app"
	"github.com/alleswebdev/mail-owl/internal/services/broker/rabbitmq"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := App.New()
	defer app.Logger.Sync()

	app.Logger.Infof("Sms-sender is running,: %#v", app.Config)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Kill)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// добавить закрытие канала

	for i := 1; i < 6; i++ {
		err := app.Broker.Subscribe(rabbitmq.EmailQueue, app.SmsHandler)

		if err != nil {
			errs <- err
		}
	}

	app.Logger.Fatalw("getting error on the errors channel", "err", <-errs)
}
