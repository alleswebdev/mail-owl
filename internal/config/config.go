package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	AppPort            int    `mapstructure:"APP_PORT"`
	AppName            string `mapstructure:"APP_NAME"`
	AppEnv             string `mapstructure:"APP_ENV"`
	AppLogPath         string `mapstructure:"APP_LOG_PATH"`
	DbHost             string `mapstructure:"DB_HOST"`
	DbPort             string `mapstructure:"DB_PORT"`
	DbName             string `mapstructure:"DB_NAME"`
	DbUser             string `mapstructure:"DB_USER"`
	DbPassword         string `mapstructure:"DB_PASSWORD"`
}

func LoadConfig(paths ...string) *Config {
	c := &Config{}
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath("../..")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/var/external/env")

	for _, val := range paths {
		viper.AddConfigPath(val)
	}

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error Config file: %s \n", err))
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to decode into struct, %v \n", err))
	}

	return c
}
