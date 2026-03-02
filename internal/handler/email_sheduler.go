package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/marisasha/email-scheduler/internal/models"
)

// @Summary Создание нового напоминания
// @Tags email-scheduler
// @ID reminder-create
// @Accept json
// @Produce json
// @Param input body models.Reminder true "Информация для напоминания"
// @Security ApiKeyAuth
// @Router /api/email-scheduler/reminder/create [post]
func (h *Handler) createReminder(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input models.Reminder

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.EmailScheduler.CreateReminder(&userId, &input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusCreated, map[string]string{
		"message": "reminder sucsuccessfully created",
	})
}

// @Summary Создание временного диапазона для напоминаний
// @Tags email-scheduler
// @ID reminder-create-range
// @Accept json
// @Produce json
// @Param input body models.RemindersWithTimeRange true "Информация для создания диапазона напоминаний"
// @Security ApiKeyAuth
// @Router /api/email-scheduler/reminder/create-range [post]
func (h *Handler) createReminderRange(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input models.RemindersWithTimeRange

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.EmailScheduler.CreateReminderRange(&userId, &input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusCreated, map[string]string{
		"message": "reminders sucsuccessfully created",
	})

}

// @Summary Просмотр напоминаний
// @Tags email-scheduler
// @ID get-reminders
// @Accept json
// @Produce json
// @Param status query string false "Статус напоминания"
// @Security ApiKeyAuth
// @Router /api/email-scheduler/reminder [get]
func (h *Handler) getReminders(c *gin.Context) {

	status := c.Query("status")

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	reminders, err := h.services.EmailScheduler.GetReminders(&userId, &status)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusCreated, map[string][]models.Reminder{
		"Data": reminders,
	})
}

// @Summary Удалить напоминание
// @Tags email-scheduler
// @ID delete-reminder
// @Accept json
// @Produce json
// @Param id path int true "Id reminder"
// @Security ApiKeyAuth
// @Router /api/email-scheduler/reminder/delete/{id} [delete]
func (h *Handler) deleteReminder(c *gin.Context) {
	reminderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadGateway, err.Error())
	}

	if err := h.services.EmailScheduler.DeleteReminder(&reminderId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusCreated, map[string]string{
		"message": "reminders sucsuccessfully deleted",
	})

}
