package email

import (
	"encoding/json"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gopkg.in/gomail.v2"
)

func RunEmailWorker(rabbitURL, queueName, smtpHost string, smtpUser string, smtpPass string, smtpPort int) error {
	// Подключаемся к RabbitMQ
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
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	logrus.Println("Email worker запущен, ждём задачи...")

	// Горрутина слушает очередь
	go func() {
		for d := range msgs {
			job := EmailJob{}
			if err := json.Unmarshal(d.Body, &job); err != nil {
				log.Println("Ошибка парсинга задачи:", err)
				continue
			}

			// Отправка письма через SMTP
			m := gomail.NewMessage()
			m.SetHeader("From", smtpUser)
			m.SetHeader("To", job.To)
			m.SetHeader("Subject", job.Subject)
			m.SetBody("text/plain", job.Body)

			dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
			if err := dialer.DialAndSend(m); err != nil {
				logrus.Println("Ошибка отправки письма:", err)
				continue
			}

			logrus.Println("Email успешно отправлен:", job.To)
		}
	}()

	// Блокируем main, чтобы горутина работала постоянно
	select {}
}
