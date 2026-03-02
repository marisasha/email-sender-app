package handler

import (
	"github.com/marisasha/email-scheduler/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	router.Use(h.loggingMiddleware)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		verifyEmail := auth.Group("/verify-email")
		{
			verifyEmail.GET("/send", h.userIdentity, h.sendEmailVerification)
			verifyEmail.GET("/check", h.checkEmailVerification)
		}
	}

	api := router.Group("/api", h.userIdentity)
	{
		emailScheduler := api.Group("/email-scheduler")
		{
			reminder := emailScheduler.Group("/reminder")
			{
				reminder.POST("/create", h.createReminder)
				reminder.POST("/create-range", h.createReminderRange)
				reminder.GET("/", h.getReminders)
				reminder.DELETE("/delete/:id", h.deleteReminder)
			}
		}
	}

	return router

}
