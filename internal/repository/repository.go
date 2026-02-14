package repository

import (
	"github.com/marisasha/email-scheduler/internal/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user *models.User) error
	GetUser(username, password string) (models.User, error)
	CreateEmailVerificationToken(userId *int, token *string) error
	CheckVerificationToken(token *string) (*models.EmailVerification, error)
	ChangeEmailVerificationStatus(userId *int) error
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthMySQL(db),
	}
}
