package main

import (
	"github.com/marisasha/email-scheduler/internal/email"
	"github.com/marisasha/email-scheduler/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logger.Init()
	if err := initConfig(); err != nil {
		logrus.Fatal(err)
	}

	err := email.RunEmailWorker(
		viper.GetString("rabbitmq_url"),
		"email_queue",
		viper.GetString("email.host"),
		viper.GetString("email.user"),
		viper.GetString("email.password"),
		viper.GetInt("email.port"),
	)
	if err != nil {
		logrus.Fatal(err)
	}
}

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	return viper.ReadInConfig()
}
