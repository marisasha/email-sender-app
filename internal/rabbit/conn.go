package rabbit

import "github.com/streadway/amqp"

type Rabbit struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbit(rabbitURL string) (*Rabbit, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Rabbit{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (r *Rabbit) DeclareQueue(queueName string) error {
	_, err := r.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	return err
}

func (r *Rabbit) Publish(queueName string, body []byte) error {
	return r.Channel.Publish(
		"",        // default exchange
		queueName, // имя очереди
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
