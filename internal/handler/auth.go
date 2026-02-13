package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marisasha/email-scheduler/internal/models"
)

// @Summary Регистрация пользователя
// @Tags auth
// @Description Создание нового пользователя
// @ID sign-up
// @Accept json
// @Produce json
// @Param input body models.User true "Данные пользователя"
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Authorization.CreateUser(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]string{
		"message": "account sucsuccessfully created",
	})

}

// @Summary Аутентификация пользователя
// @Tags auth
// @Description Проверка прав пользователя
// @ID sign-in
// @Accept json
// @Produce json
// @Param input body models.UserSignInRequest true "Данные пользователя"
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input models.UserSignInRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

type sendEmailRequest struct {
	Email string `json:"email"`
}

// @Summary Подтверждение почты пользователя
// @Tags auth
// @Description Проверка почты пользователя
// @ID verify
// @Accept json
// @Produce json
// @Param input body sendEmailRequest true "почта пользователя"
// @Security ApiKeyAuth
// @Router /auth/verify-token [post]
func (h *Handler) sendEmailVerification(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input sendEmailRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Authorization.SendEmailVerification(&userId, &input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]string{
		"message": "message with email verification sucsuccessfully sent",
	})
}
