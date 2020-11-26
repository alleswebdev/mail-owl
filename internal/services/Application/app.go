package Application

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
	"vcs.torbor.ru/notizer/workers/spectacler/internal/config"
	"vcs.torbor.ru/notizer/workers/spectacler/internal/log"
	"vcs.torbor.ru/notizer/workers/spectacler/internal/services/NoticeBuilder"
	"vcs.torbor.ru/notizer/workers/spectacler/internal/storage"
)

type App struct {
	Config  config.Config
	Logger  zap.SugaredLogger
	Storage storage.DBStorage
	Builder NoticeBuilder.Builder
	Metrics *Metrics
}

type Metrics struct {
	PgGauge prometheus.Gauge
	PgStat  prometheus.Gauge
	Up      prometheus.Gauge
}

func New() *App {
	cfg := config.LoadConfig()

	var logger zap.SugaredLogger

	if strings.Contains(strings.ToLower(cfg.AppEnv), "dev") {
		logger = log.BuildDevLogger(cfg.AppLogPath)
	} else {
		logger = log.BuildProdLogger(cfg.AppLogPath)
	}

	db, err := connDb(cfg, logger)

	if err != nil {
		logger.Fatal(err)
	}

	s := storage.NewStorage(*db)

	tplStorage, err := storage.NewTemplateStorage(*cfg)

	if err != nil {
		logger.Fatal(err)
	}

	return &App{
		Config:  *cfg,
		Logger:  logger,
		Storage: s,
		Builder: *NoticeBuilder.New(tplStorage, cfg),
		Metrics: initMetrics(cfg),
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

func initMetrics(cfg *config.Config) *Metrics {
	m := &Metrics{}
	m.PgGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "app_db_connection_check",
		Help:        "database check",
		ConstLabels: prometheus.Labels{"db_type": "pgsql", "db_name": cfg.DbPort, "db_host": cfg.DbHost},
	})

	m.Up = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "up",
		Help:        "app info, this is NOT A HEALZ-CHECK!!!!!!!",
		ConstLabels: prometheus.Labels{"environment": cfg.AppEnv, "name": cfg.AppName, "port": strconv.Itoa(cfg.AppPort)},
	})

	prometheus.MustRegister(m.PgGauge, m.Up)

	return m
}
