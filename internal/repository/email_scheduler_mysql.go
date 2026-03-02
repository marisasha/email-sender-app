package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/marisasha/email-scheduler/internal/models"
)

type EmailSchedulerMySQL struct {
	db *sqlx.DB
}

func NewEmailSchedulerMySQL(db *sqlx.DB) *EmailSchedulerMySQL {
	return &EmailSchedulerMySQL{db: db}
}

func (r *EmailSchedulerMySQL) CreateReminder(userId *int, input *models.Reminder) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id,subject,text,send_at) VALUES(?,?,?,?)", emailReminderTable)
	_, err := r.db.Exec(query, *userId, input.Subject, input.Text, input.SentAt)
	return err
}

func (r *EmailSchedulerMySQL) GetPendingReminders() ([]models.Reminder, error) {

	var reminders []models.Reminder

	query := fmt.Sprintf(`
		SELECT *
			FROM %s
			WHERE status = 'pending'
			AND send_at <= NOW()
			ORDER BY send_at
	`, emailReminderTable)

	err := r.db.Select(&reminders, query)
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

func (r *EmailSchedulerMySQL) UpdateReminderStatus(id int, status string) error {
	query := fmt.Sprintf("UPDATE %s SET status = ? WHERE id = ?", emailReminderTable)
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *EmailSchedulerMySQL) GetEmail(userId *int) (string, error) {

	var email string

	query := fmt.Sprintf("SELECT email FROM %s WHERE id=?", userTable)
	row := r.db.QueryRow(query, userId)
	if err := row.Scan(&email); err != nil {
		return "", err
	}
	return email, nil
}

func (r *EmailSchedulerMySQL) CreateReminders(reminders []models.Reminder) error {

	query := fmt.Sprintf("INSERT INTO %s (user_id, subject, text, send_at) VALUES ", emailReminderTable)

	var args []interface{}

	for i, reminder := range reminders {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?, ?)"

		args = append(args,
			reminder.UserId,
			reminder.Subject,
			reminder.Text,
			reminder.SentAt,
		)
	}

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *EmailSchedulerMySQL) GetReminders(userId *int, status *string) ([]models.Reminder, error) {
	var reminders []models.Reminder

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = ? AND status=?", emailReminderTable)

	if err := r.db.Select(&reminders, query, *userId, &status); err != nil {
		return nil, err
	}

	return reminders, nil
}

func (r *EmailSchedulerMySQL) DeleteReminder(reminderId *int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=?", emailReminderTable)
	_, err := r.db.Exec(query, *reminderId)
	return err
}
