package service

import (
	emailservice "github.com/marisasha/email-scheduler/internal/email"
	"github.com/marisasha/email-scheduler/internal/models"
	"github.com/marisasha/email-scheduler/internal/repository"
)

type Authorization interface {
	CreateUser(user *models.User) error
	GenerateToken(username, password *string) (string, error)
	ParseToken(token *string) (int, error)
	SendEmailVerification(userId *int) error
	CheckEmailVerification(token *string) error
}

type EmailScheduler interface {
	CreateReminder(userId *int, input *models.Reminder) error
	CreateReminderRange(userId *int, input *models.RemindersWithTimeRange) error
	GetReminders(userId *int, status *string) ([]models.Reminder, error)
	DeleteReminder(reminderId *int) error
	StartScheduler()
}
type Service struct {
	Authorization
	EmailScheduler
}

func NewService(repos *repository.Repository, emailQueue *emailservice.EmailService) *Service {
	return &Service{
		Authorization:  NewAuthService(repos.Authorization, emailQueue.Publisher),
		EmailScheduler: NewEmailSchedulerService(repos.EmailScheduler, emailQueue.Publisher),
	}
}
