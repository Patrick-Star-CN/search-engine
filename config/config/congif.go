package config

import (
	"github.com/spf13/viper"
	"log"
)

var Config = viper.New()

func init() {
	Config.SetConfigName("config")
	Config.SetConfigType("yaml")
	Config.AddConfigPath(".")
	Config.WatchConfig()
	err := Config.ReadInConfig()
	if err != nil {
		log.Fatal("Config not find", err)
	}
}
