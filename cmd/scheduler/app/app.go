package App

import (
	"context"
	"fmt"
	"github.com/alleswebdev/mail-owl/internal/config"
	"github.com/alleswebdev/mail-owl/internal/log"
	"github.com/alleswebdev/mail-owl/internal/models"
	"github.com/alleswebdev/mail-owl/internal/services/broker"
	"github.com/alleswebdev/mail-owl/internal/services/broker/rabbitmq"
	"github.com/alleswebdev/mail-owl/internal/storage"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"time"
)

type App struct {
	Config  config.Config
	Logger  zap.SugaredLogger
	Storage storage.DBStorage
	Broker  broker.Broker
}

func New() *App {
	cfg := config.LoadConfig()

	var logger zap.SugaredLogger

	logger = log.BuildMainLogger(cfg.AppLogPath)

	db, err := connDb(cfg, logger)

	if err != nil {
		logger.Fatal(err)
	}

	s := storage.NewStorage(*db)

	if err != nil {
		logger.Fatal(err)
	}

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
		Config:  *cfg,
		Logger:  logger,
		Storage: s,
		Broker:  b,
	}
}

func connDb(cfg *config.Config, logger zap.SugaredLogger) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s pool_max_conns=5",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName)

	poolConfig, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, fmt.Errorf("unable to parse DATABASE_URL %s", err)
	}

	poolConfig.MaxConnLifetime = 15 * time.Minute
	poolConfig.ConnConfig.PreferSimpleProtocol = true
	poolConfig.ConnConfig.Logger = zapadapter.NewLogger(logger.Desugar())

	dbpool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)

	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	return dbpool, nil
}

type PublishState struct {
	Notice models.SchedulerNotice
	State  models.SchedulerState
	Error  error
	Queue  string
}

func (app *App) PublishState(event PublishState) {
	event.Notice.State = event.State

	if event.Error != nil {
		event.Notice.Error = event.Error.Error()
	}

	noticeEncode, err := event.Notice.MarshalJSON()

	if err != nil {
		app.Logger.Errorf("error on marshal json for id %d, state:%s, err: %s", event.Notice.Id, event.Notice.State, err)
		return
	}

	if event.Queue == "" {
		event.Queue = rabbitmq.SchedulerQueue
	}

	err = app.Broker.Publish(broker.Message{
		Body:    noticeEncode,
		Headers: nil,
	}, event.Queue)

	if err != nil {
		app.Logger.Errorf("error on publish notice state for id %d, state:%s, err: %s", event.Notice.Id, event.Notice.State, err)
		return
	}
}
