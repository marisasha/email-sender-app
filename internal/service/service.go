package service

import (
	"github.com/marisasha/email-scheduler/internal/email"
	"github.com/marisasha/email-scheduler/internal/models"
	"github.com/marisasha/email-scheduler/internal/repository"
)

type Authorization interface {
	CreateUser(user *models.User) error
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	SendEmailVerification(userId *int, email *string) error
	CheckEmailVerification(token *string) error
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository, emailQueue *email.Email) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, emailQueue.EmailRepository),
	}
}
