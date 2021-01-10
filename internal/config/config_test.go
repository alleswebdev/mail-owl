package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	b := []byte(`
				DB_PORT=5432
				DB_NAME=mailowl
				RABBIT_USER=rabbit-test
			`)
	path := ".env"
	err := ioutil.WriteFile(path, b, 0644)
	if err != nil {
		t.Error(err)
	}

	cfg := LoadConfig()

	if cfg.RabbitUser != "rabbit-test" {
		t.Error("RabbitUser is not valid")
	}

	if cfg.DbName != "mailowl" {
		t.Error("DbName is not valid")
	}

	if cfg.DbPort != "5432" {
		t.Error("DbPort is not valid")
	}
	_ = os.Remove(path)
}
