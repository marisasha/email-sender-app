package repository

import (
	"github.com/marisasha/email-scheduler/internal/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user *models.User) error
	GetUser(username, password *string) (models.User, error)
	CreateEmailVerificationToken(userId *int, token *string) error
	GetUserEmail(userId *int) (string, error)
	CheckVerificationToken(token *string) (*models.EmailVerification, error)
	ChangeEmailVerificationStatus(userId *int) error
}

type EmailScheduler interface {
	CreateReminder(userId *int, input *models.Reminder) error
	CreateReminders(reminders []models.Reminder) error
	GetPendingReminders() ([]models.Reminder, error)
	UpdateReminderStatus(id int, status string) error
	GetEmail(userId *int) (string, error)
	GetReminders(userId *int, status *string) ([]models.Reminder, error)
	DeleteReminder(reminderId *int) error
}

type Repository struct {
	Authorization
	EmailScheduler
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization:  NewAuthMySQL(db),
		EmailScheduler: NewEmailSchedulerMySQL(db),
	}
}
