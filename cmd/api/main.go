package main

import (
	"context"
	"log"

	"github.com/malcolmkzh/study-notifier/internal/modules/healthcheck"
	"github.com/malcolmkzh/study-notifier/internal/modules/notes"
	"github.com/malcolmkzh/study-notifier/internal/modules/questions"
	"github.com/malcolmkzh/study-notifier/internal/modules/reminder"
	"github.com/malcolmkzh/study-notifier/internal/utilities/config"
	"github.com/malcolmkzh/study-notifier/internal/utilities/db"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpclient"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpserver"
	"github.com/malcolmkzh/study-notifier/internal/utilities/llm"
	"github.com/malcolmkzh/study-notifier/internal/utilities/notification"
	"github.com/malcolmkzh/study-notifier/internal/utilities/scheduler"
)

type AppDependencies struct {
	DB           db.Utility
	HTTPServer   httpserver.Utility
	LLM          llm.Utility
	Scheduler    scheduler.Utility
	JobRepo      scheduler.JobRepository
	Notification notification.Utility
}

func main() {
	// startup context
	ctx := context.Background()

	//Init utilities
	configUtility, err := config.NewConfigUtility()
	if err != nil {
		log.Fatal("Failed to initialize configuration utility: ", err)
	}

	dbUtility, err := db.NewDbUtility(configUtility)
	if err != nil {
		log.Fatal("Failed to initialize database utility: ", err)
	}

	httpServerUtility := httpserver.NewHttpServerUtility(configUtility)
	httpClientUtility := httpclient.NewHTTPClientUtility()
	llmUtility, err := llm.NewLLMUtility(configUtility, httpClientUtility)
	if err != nil {
		log.Fatal("Failed to initialize llm utility: ", err)
	}
	jobRepository, err := scheduler.NewJobRepository(dbUtility)
	if err != nil {
		log.Fatal("Failed to initialize scheduler job repository: ", err)
	}
	schedulerUtility, err := scheduler.NewUtility(jobRepository)
	if err != nil {
		log.Fatal("Failed to initialize scheduler utility: ", err)
	}
	notificationUtility := notification.NewNotificationUtility()

	appDependencies := AppDependencies{
		DB:           dbUtility,
		HTTPServer:   httpServerUtility,
		LLM:          llmUtility,
		Scheduler:    schedulerUtility,
		JobRepo:      jobRepository,
		Notification: notificationUtility,
	}

	//Init Modules
	_, err = notes.New(ctx, notes.Dependencies{
		DB:         appDependencies.DB,
		HTTPServer: appDependencies.HTTPServer,
	})
	if err != nil {
		log.Fatal("Failed to initialize notes module: ", err)
	}

	_, err = questions.New(ctx, questions.Dependencies{
		DB:         appDependencies.DB,
		LLM:        appDependencies.LLM,
		HTTPServer: appDependencies.HTTPServer,
	})
	if err != nil {
		log.Fatal("Failed to initialize questions module: ", err)
	}

	_, err = healthcheck.New(ctx, healthcheck.Dependencies{
		HTTPServer: appDependencies.HTTPServer,
	})
	if err != nil {
		log.Fatal("Failed to initialize healthcheck module: ", err)
	}

	_, err = reminder.New(ctx, reminder.Dependencies{
		DB:           appDependencies.DB,
		Scheduler:    appDependencies.Scheduler,
		JobRepo:      appDependencies.JobRepo,
		Notification: appDependencies.Notification,
	})
	if err != nil {
		log.Fatal("Failed to initialize reminder module: ", err)
	}

	err = appDependencies.Scheduler.Start(ctx)
	if err != nil {
		log.Fatal("Failed to start scheduler utility: ", err)
	}

	// Start HTTP server
	err = httpServerUtility.Serve(context.Background())
	if err != nil {
		log.Fatal("Failed to start HTTP server: ", err)
	}

	select {}
}
