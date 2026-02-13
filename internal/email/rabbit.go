package email

import "github.com/streadway/amqp"

type EmailRabbit struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewEmailRepository(rabbitURL, queueName string) (*EmailRabbit, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &EmailRabbit{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}
