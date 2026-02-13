package email

type EmailRepository interface {
	PublishEmail(job EmailJob) error
	Close()
}

type Email struct {
	EmailRepository
}

func NewEmail(emailQueue *EmailRabbit) *Email {
	return &Email{
		EmailRepository: NewAuthEmail(emailQueue),
	}
}
