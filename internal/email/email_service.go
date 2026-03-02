package emailservice

import (
	"github.com/marisasha/email-scheduler/internal/models"
	"github.com/marisasha/email-scheduler/internal/rabbit"
)

type Publisher interface {
	PublishEmail(job models.EmailJob, queueName string) error
	Close()
}

type EmailService struct {
	Publisher
}

func NewEmailService(rabbit *rabbit.Rabbit) *EmailService {
	return &EmailService{
		Publisher: NewEmailPublisher(rabbit),
	}
}
