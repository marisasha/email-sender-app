package main

import (
	emailservice "github.com/marisasha/email-scheduler/internal/email"
	"github.com/marisasha/email-scheduler/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logger.Init()
	if err := initConfig(); err != nil {
		logrus.Fatal(err)
	}

	rabbitmqUrl := viper.GetString("rabbitmq_url")
	host := viper.GetString("email.host")
	user := viper.GetString("email.user")
	password := viper.GetString("email.password")
	port := viper.GetInt("email.port")

	go func() {
		if err := emailservice.RunEmailWorker(
			rabbitmqUrl,
			"email_reminder",
			host,
			user,
			password,
			port,
		); err != nil {
			logrus.Fatal(err)
		}
	}()

	go func() {
		if err := emailservice.RunEmailWorker(
			rabbitmqUrl,
			"email_verification",
			host,
			user,
			password,
			port,
		); err != nil {
			logrus.Fatal(err)
		}
	}()

	select {}
}

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	return viper.ReadInConfig()
}
