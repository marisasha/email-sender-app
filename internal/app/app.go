package app

import (
	"context"

	"github.com/marisasha/email-scheduler/internal/email"
	"github.com/marisasha/email-scheduler/internal/handler"
	"github.com/marisasha/email-scheduler/internal/repository"
	"github.com/marisasha/email-scheduler/internal/service"
	httpserver "github.com/marisasha/email-scheduler/internal/transport/http"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type App struct {
	server     *httpserver.Server
	handlers   *handler.Handler
	db         *sqlx.DB
	emailQueue *email.EmailRabbit
}

func NewApp(cfg repository.Config) (*App, error) {
	db, err := repository.NewMySQLDB(cfg)
	if err != nil {
		return nil, err
	}

	queue, err := email.NewEmailRepository(viper.GetString("rabbitmq_url"), "email_queue")
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepository(db)
	emailQueue := email.NewEmail(queue)
	services := service.NewService(repos, emailQueue)
	handlers := handler.NewHandler(services)

	server := new(httpserver.Server)

	return &App{
		server:     server,
		handlers:   handlers,
		db:         db,
		emailQueue: queue,
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
