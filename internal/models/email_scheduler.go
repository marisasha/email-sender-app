package models

import "time"

type Reminder struct {
	Id      int       `json:"id" db:"id"`
	UserId  int       `json:"-" db:"user_id"`
	Subject string    `json:"subject" db:"subject" binding:"required" `
	Text    string    `json:"text" db:"text" binding:"required"`
	SentAt  time.Time `json:"send_at" db:"send_at" binding:"required"`
	Status  string    `json:"status" db:"status" default:"pending"`
}

type RemindersWithTimeRange struct {
	Id              int       `json:"-" `
	UserId          int       `json:"-" `
	Subject         string    `json:"subject" binding:"required" `
	Text            string    `json:"text" binding:"required"`
	RangeStart      time.Time `json:"range_start"  binding:"required"`
	RangeEnd        time.Time `json:"range_end"  binding:"required"`
	RepeatCondition string    `json:"repeat_condition" binding:"required,oneof=year month week day"`
	Condition       string    `json:"condition"` //
}

type EmailJob struct {
	To      string
	Subject string
	Body    string
}
