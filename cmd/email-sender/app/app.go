package App

import (
	"fmt"
	"github.com/alleswebdev/mail-owl/internal/config"
	"github.com/alleswebdev/mail-owl/internal/log"
	"github.com/alleswebdev/mail-owl/internal/models"
	"github.com/alleswebdev/mail-owl/internal/services/broker"
	"github.com/alleswebdev/mail-owl/internal/services/broker/rabbitmq"
	"github.com/alleswebdev/mail-owl/internal/services/mail"
	"go.uber.org/zap"
)

type App struct {
	Config config.Config
	Logger zap.SugaredLogger
	Broker broker.Broker
	Mailer mail.Mailer
}

func New() *App {
	cfg := config.LoadConfig()

	var logger zap.SugaredLogger

	logger = log.BuildMainLogger(cfg.AppLogPath)

	uri := fmt.Sprintf("amqp://%s:%s@%s:%s",
		cfg.RabbitUser,
		cfg.RabbitPassword,
		cfg.RabbitHost,
		cfg.RabbitPort,
	)

	b, err := rabbitmq.NewRabbit(uri, "mailowl-default", &logger)

	if err != nil {
		logger.Fatal(err)
	}

	mailer := mail.NewMailer(*cfg)

	return &App{
		Config: *cfg,
		Logger: logger,
		Broker: b,
		Mailer: *mailer,
	}
}

func (app *App) PublishState(notice models.SchedulerNotice, state models.SchedulerState, error error) {
	notice.State = state

	if error != nil {
		notice.Error = error.Error()
	}

	noticeEncode, err := notice.MarshalJSON()

	if err != nil {
		app.Logger.Errorf("error on marshal json for id %d, state:%s, err: %s", notice.Id, notice.State, err)
		return
	}

	err = app.Broker.Publish(broker.Message{
		Body:    noticeEncode,
		Headers: nil,
	}, rabbitmq.SchedulerQueue)

	if err != nil {
		app.Logger.Errorf("error on publish notice state for id %d, state:%s, err: %s", notice.Id, notice.State, err)
		return
	}
}
