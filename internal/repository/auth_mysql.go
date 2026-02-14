package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/marisasha/email-scheduler/internal/models"
)

type AuthMySQL struct {
	db *sqlx.DB
}

func NewAuthMySQL(db *sqlx.DB) *AuthMySQL {
	return &AuthMySQL{db: db}
}

func (r *AuthMySQL) CreateUser(user *models.User) error {
	query := fmt.Sprintf("INSERT INTO %s (email,password_hash,first_name,last_name) VALUES (?,?,?,?)", userTable)
	_, err := r.db.Exec(query, user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return err
	}
	return nil

}

func (r *AuthMySQL) GetUser(username, password string) (models.User, error) {

	var user models.User

	query := fmt.Sprintf("SELECT id FROM %s WHERE email=? AND password_hash=?", userTable)
	err := r.db.Get(&user, query, username, password)

	return user, err

}

func (r *AuthMySQL) CreateEmailVerificationToken(userId *int, token *string) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id,token) VALUES (?,?)", verificationTokenTable)
	_, err := r.db.Exec(query, *userId, *token)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthMySQL) CheckVerificationToken(token *string) (*models.EmailVerification, error) {

	var emailVerification models.EmailVerification
	query := fmt.Sprintf("SELECT * FROM %s WHERE token=? ", verificationTokenTable)
	err := r.db.Get(&emailVerification, query, *token)
	if err != nil {
		return nil, err
	}

	return &emailVerification, nil
}

func (r *AuthMySQL) ChangeEmailVerificationStatus(userId *int) error {
	query := fmt.Sprintf("UPDATE %s SET email_verified=1 WHERE id=?", userTable)
	_, err := r.db.Exec(query, *userId)
	if err != nil {
		return err
	}
	return nil

}
