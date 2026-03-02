package service

import (
	"time"

	"github.com/marisasha/email-scheduler/internal/models"
	"github.com/sirupsen/logrus"
)

func (s *EmailSchedulerService) StartScheduler() {
	logrus.Info("Scheduler is starting !")
	go func() {
		for {
			reminders, err := s.repos.GetPendingReminders()
			if err != nil {
				logrus.Error(err)
				time.Sleep(10 * time.Second)
				continue
			}

			if len(reminders) == 0 {
				logrus.Info("Нет необработанных напоминаний")
				time.Sleep(30 * time.Second)
				continue
			}

			for _, r := range reminders {
				userEmail, err := s.repos.GetEmail(&r.UserId)
				if err != nil {
					logrus.Error("Ошибка получения напоминаний:", err)
					continue
				}
				job := models.EmailJob{
					To:      userEmail,
					Subject: r.Subject,
					Body:    r.Text,
				}

				err = s.emailQueue.PublishEmail(job, "email_reminder")
				if err != nil {
					logrus.Error("Ошибка публикации в RabbitMQ:", err)
					continue
				}

				err = s.repos.UpdateReminderStatus(r.Id, "sent")
				if err != nil {
					logrus.Error("Ошибка обновления статуса:", err)
					continue
				}
				logrus.Info("!!! Напоминание отправленно в очередь !!!")
				time.Sleep(30 * time.Second)
			}
		}
	}()
}
