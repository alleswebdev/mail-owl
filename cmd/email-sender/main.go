package main

import (
	"fmt"
	"github.com/alleswebdev/mail-owl/cmd/email-sender/app"
	"github.com/alleswebdev/mail-owl/internal/services/broker/rabbitmq"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := App.New()
	defer app.Logger.Sync()

	app.Logger.Infof("Email-sender is running,: %#v", app.Config)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Kill)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// добавить закрытие канала

	for i := 1; i < 6; i++ {
		err := app.Broker.Subscribe(rabbitmq.EmailQueue, app.EmailHandler)

		if err != nil {
			errs <- err
		}
	}

	// close idle smtp connection, it need a mux on mailer because data race
	//go func() {
	//	for {
	//		select {
	//		case <-time.After(10 * time.Second):
	//			if app.Mailer.Open {
	//				app.Logger.Info("smtp connection is idle 10s, closing the connection")
	//
	//				if err := app.Mailer.Shutdown(); err != nil {
	//					app.Logger.Error(err)
	//					app.Mailer.Open = false
	//					continue
	//				}
	//
	//				app.Mailer.Open = false
	//				log.Info("smtp connection was close")
	//			}
	//		}
	//	}
	//}()

	app.Logger.Fatalw("getting error on the errors channel", "err", <-errs)
}
