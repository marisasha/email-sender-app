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
