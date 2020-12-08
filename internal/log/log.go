package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func BuildMainLogger(path string) zap.SugaredLogger {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"stderr", path}
	cfg.ErrorOutputPaths = []string{"stderr", path}
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.TimeKey = ""

	return BuildLogger(path, cfg)
}
