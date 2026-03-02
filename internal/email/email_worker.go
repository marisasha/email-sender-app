package emailservice

import (
	"encoding/json"
	"log"

	"github.com/marisasha/email-scheduler/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gopkg.in/gomail.v2"
)

func RunEmailWorker(rabbitURL, queueName, smtpHost string, smtpUser string, smtpPass string, smtpPort int) error {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	logrus.Printf("Worker %s запущен, ждём задачи...", queueName)

	go func() {
		for d := range msgs {
			job := models.EmailJob{}
			if err := json.Unmarshal(d.Body, &job); err != nil {
				log.Println("Ошибка парсинга задачи:", err)
				continue
			}

			m := gomail.NewMessage()
			m.SetHeader("From", smtpUser)
			m.SetHeader("To", job.To)
			m.SetHeader("Subject", job.Subject)
			m.SetBody("text/plain", job.Body)

			dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
			if err := dialer.DialAndSend(m); err != nil {
				logrus.Printf("%s:Ошибка отправки письма:%s", queueName, err)
				continue
			}

			logrus.Printf("Воркер %s успешно отправил email по адресу: %s", queueName, job.To)
		}
	}()

	select {}
}
