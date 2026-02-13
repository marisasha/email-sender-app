package main

import (
	"github.com/marisasha/email-scheduler/internal/email"
	"github.com/marisasha/email-scheduler/internal/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.Init()
	err := email.RunEmailWorker(
		"amqp://guest:guest@127.0.0.1:5672/", // RabbitMQ URL
		"email_queue",                        // Имя очереди
		"smtp.mail.ru",                       // SMTP хост
		"marisasham987@mail.ru",              // SMTP пользователь
		"P30V3dsqeE5ScYPtF4cQ",               // SMTP пароль
		465,                                  // SMTP порт
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
