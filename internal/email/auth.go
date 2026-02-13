package email

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type AuthEmail struct {
	emailQueue *EmailRabbit
}

func NewAuthEmail(emailQueue *EmailRabbit) *AuthEmail {
	return &AuthEmail{emailQueue: emailQueue}
}

type EmailJob struct {
	To      string
	Subject string
	Body    string
}

func (r *AuthEmail) PublishEmail(job EmailJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return r.emailQueue.channel.Publish(
		"",                 // exchange
		r.emailQueue.queue, // routing key = имя очереди
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Close закрывает соединение с RabbitMQ
func (r *AuthEmail) Close() {
	if r.emailQueue.channel != nil {
		r.emailQueue.channel.Close()
	}
	if r.emailQueue.conn != nil {
		r.emailQueue.conn.Close()
	}
}
