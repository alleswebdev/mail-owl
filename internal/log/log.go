package log

import (
	"go.uber.org/zap"
	"os"
	"path"
)

func BuildLogger(p string, cfg zap.Config) zap.SugaredLogger {
	err := os.MkdirAll(path.Dir(p), 0755)

	if err != nil {
		panic(err)
	}

	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}

	return *logger.Sugar()
}

