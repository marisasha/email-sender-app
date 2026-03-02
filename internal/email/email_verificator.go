package emailservice

import (
	"encoding/json"

	"github.com/marisasha/email-scheduler/internal/models"
	"github.com/marisasha/email-scheduler/internal/rabbit"
)

type EmailPublisher struct {
	rabbit *rabbit.Rabbit
}

func NewEmailPublisher(rabbit *rabbit.Rabbit) *EmailPublisher {
	return &EmailPublisher{rabbit: rabbit}
}

func (p *EmailPublisher) PublishEmail(job models.EmailJob, queueName string) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return p.rabbit.Publish(queueName, body)
}
func (p *EmailPublisher) Close() {
	if p.rabbit.Channel != nil {
		p.rabbit.Channel.Close()
	}
	if p.rabbit.Conn != nil {
		p.rabbit.Conn.Close()
	}
}
