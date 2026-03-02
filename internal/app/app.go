package app

import (
	"context"

	emailservice "github.com/marisasha/email-scheduler/internal/email"
	"github.com/marisasha/email-scheduler/internal/handler"
	"github.com/marisasha/email-scheduler/internal/rabbit"
	"github.com/marisasha/email-scheduler/internal/repository"
	"github.com/marisasha/email-scheduler/internal/service"
	httpserver "github.com/marisasha/email-scheduler/internal/transport/http"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type App struct {
	server   *httpserver.Server
	handlers *handler.Handler
	db       *sqlx.DB
}

func NewApp(cfg repository.Config) (*App, error) {
	db, err := repository.NewMySQLDB(cfg)
	if err != nil {
		return nil, err
	}

	rabbit, err := rabbit.NewRabbit(viper.GetString("rabbitmq_url"))
	if err != nil {
		return nil, err
	}
	rabbit.DeclareQueue("email_verification")
	rabbit.DeclareQueue("email_reminder")

	repos := repository.NewRepository(db)
	emailQueue := emailservice.NewEmailService(rabbit)
	services := service.NewService(repos, emailQueue)
	handlers := handler.NewHandler(services)

	server := new(httpserver.Server)
	services.EmailScheduler.StartScheduler()

	return &App{
		server:   server,
		handlers: handlers,
		db:       db,
	}, nil
}

func (a *App) Run(port string) error {
	return a.server.Run(port, a.handlers.InitRoutes())
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}
	return a.db.Close()
}
