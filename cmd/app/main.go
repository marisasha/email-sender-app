package main

import (
	_ "github.com/marisasha/email-scheduler/docs"
	"github.com/marisasha/email-scheduler/internal/app"
	"github.com/marisasha/email-scheduler/internal/config"
	_ "github.com/marisasha/email-scheduler/internal/docs"
	"github.com/marisasha/email-scheduler/internal/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("cannot load config: %s", err)
	}

	application, err := app.NewApp(cfg.DB)
	if err != nil {
		logrus.Fatalf("failed to initialize app: %s", err.Error())
	}

	app.RunWithGracefulShutdown(application, cfg.AppPort)

}
