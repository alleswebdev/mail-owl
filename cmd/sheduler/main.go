package main

import (
	"fmt"
	"github.com/alleswebdev/mail-owl/cmd/sheduler/app"
	"github.com/alleswebdev/mail-owl/internal/services/broker"
	"github.com/alleswebdev/mail-owl/internal/services/broker/rabbitmq"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := App.New()
	defer app.Storage.Db.Close()
	defer app.Logger.Sync()

	app.Logger.Infof("Sheduler is running,: %#v", app.Config)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Kill)
		errs <- fmt.Errorf("%s", <-c)
	}()

	uri := fmt.Sprintf("amqp://%s:%s@%s:%s",
		app.Config.RabbitUser,
		app.Config.RabbitPassword,
		app.Config.RabbitHost,
		app.Config.RabbitPort,
	)

	rabbit, _ := rabbitmq.NewRabbit(uri, "mailowl-default", &app.Logger)

	_, err := rabbit.Subscribe(rabbitmq.SchedulerMain, rabbitmq.ProducerHandler)

	if nil != err {
		fmt.Println(err)
	}

	err = rabbit.Publish(broker.Message{
		Body: []byte("qwer"),
	}, rabbitmq.SchedulerMain)

	if nil != err {
		fmt.Println(err)
	}

	app.Logger.Fatalw("getting error on the errors channel", "err", <-errs)
}
