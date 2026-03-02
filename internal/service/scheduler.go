package service

import (
	"errors"

	emailservice "github.com/marisasha/email-scheduler/internal/email"
	"github.com/marisasha/email-scheduler/internal/models"
	"github.com/marisasha/email-scheduler/internal/repository"
)

type EmailSchedulerService struct {
	repos      repository.EmailScheduler
	emailQueue emailservice.Publisher
}

func NewEmailSchedulerService(repos repository.EmailScheduler, emailQueue emailservice.Publisher) *EmailSchedulerService {
	return &EmailSchedulerService{
		repos:      repos,
		emailQueue: emailQueue,
	}
}

func (s *EmailSchedulerService) CreateReminder(userId *int, input *models.Reminder) error {
	return s.repos.CreateReminder(userId, input)

}

func (s *EmailSchedulerService) CreateReminderRange(userId *int, input *models.RemindersWithTimeRange) error {
	var remindersForSend []models.Reminder

	if input.RepeatCondition == "day" {
		rangeStart := input.RangeStart
		rangeEnd := input.RangeEnd
		for true {
			if !rangeStart.After(rangeEnd) {
				remindersForSend = append(remindersForSend, models.Reminder{
					Subject: input.Subject,
					Text:    input.Text,
					SentAt:  rangeStart,
				})
			} else {
				break
			}

			switch input.Condition {
			case "every_day":
				rangeStart = rangeStart.AddDate(0, 0, 1)
			case "every_other_day":
				rangeStart = rangeStart.AddDate(0, 0, 2)
			default:
				return errors.New("invalid condition")
			}
		}
	} else {
		rangeStart := input.RangeStart
		rangeEnd := input.RangeEnd

		for true {
			if !rangeStart.After(rangeEnd) {
				remindersForSend = append(remindersForSend, models.Reminder{
					UserId:  *userId,
					Subject: input.Subject,
					Text:    input.Text,
					SentAt:  rangeStart,
				})
			} else {
				break
			}
			switch input.RepeatCondition {
			case "year":
				rangeStart = rangeStart.AddDate(1, 0, 0)
			case "month":
				rangeStart = rangeStart.AddDate(0, 1, 0)
			case "week":
				rangeStart = rangeStart.AddDate(0, 0, 7)
			default:
				return errors.New("invalid repeat condition")
			}
		}
	}

	if len(remindersForSend) == 0 {
		return errors.New("no reminders were created in the specified range")
	}

	return s.repos.CreateReminders(remindersForSend)
}

func (s *EmailSchedulerService) GetReminders(userId *int, status *string) ([]models.Reminder, error) {
	return s.repos.GetReminders(userId, status)
}

func (s *EmailSchedulerService) DeleteReminder(reminderId *int) error {
	return s.repos.DeleteReminder(reminderId)
}
