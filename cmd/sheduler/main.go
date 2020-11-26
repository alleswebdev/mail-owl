package main

import (
	"fmt"
	"github.com/alleswebdev/mail-owl/cmd/sheduler/app"
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

	go func() {
		//
	}()

	app.Logger.Fatalw("getting error on the errors channel", "err", <-errs)
}
