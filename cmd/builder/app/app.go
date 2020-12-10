package App

import (
	"fmt"
	"github.com/alleswebdev/mail-owl/internal/config"
	"github.com/alleswebdev/mail-owl/internal/log"
	"github.com/alleswebdev/mail-owl/internal/services/broker"
	"github.com/alleswebdev/mail-owl/internal/services/broker/rabbitmq"
	"go.uber.org/zap"
)

type App struct {
	Config config.Config
	Logger zap.SugaredLogger
	Broker broker.Broker
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

	return &App{
		Config: *cfg,
		Logger: logger,
		Broker: b,
	}
}
