package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestBuildMainLogger(t *testing.T) {
	path := "./test.log"
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{path}
	cfg.ErrorOutputPaths = []string{path}
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.TimeKey = ""

	logger := BuildLogger(path, cfg)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error(err)
	}
	testText := "logger work"
	logger.Info(testText)
	data, err := ioutil.ReadFile(path)

	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(data), testText) {
		t.Error(err)
	}

	_ = os.Remove(path)
}
