package reminder

import (
	"context"
	"errors"

	questionrepository "github.com/malcolmkzh/study-notifier/internal/modules/questions/repository"
	"github.com/malcolmkzh/study-notifier/internal/modules/reminder/controller"
	"github.com/malcolmkzh/study-notifier/internal/modules/reminder/repository"
	"github.com/malcolmkzh/study-notifier/internal/modules/reminder/service"
	userrepository "github.com/malcolmkzh/study-notifier/internal/modules/user/repository"
	"github.com/malcolmkzh/study-notifier/internal/utilities/db"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpserver"
	"github.com/malcolmkzh/study-notifier/internal/utilities/notification"
	"github.com/malcolmkzh/study-notifier/internal/utilities/scheduler"
)

type Dependencies struct {
	DB           db.Utility
	HTTPServer   httpserver.Utility
	Scheduler    scheduler.Utility
	JobRepo      scheduler.JobRepository
	Notification notification.Utility
}

type Module struct {
	Repository repository.Utility
	Service    service.Service
	Controller *controller.Implementation
}

func New(ctx context.Context, dependencies Dependencies) (*Module, error) {
	_ = ctx

	if dependencies.DB == nil {
		return nil, errors.New("db dependency is required")
	}
	if dependencies.HTTPServer == nil {
		return nil, errors.New("http server dependency is required")
	}
	if dependencies.Scheduler == nil {
		return nil, errors.New("scheduler dependency is required")
	}
	if dependencies.JobRepo == nil {
		return nil, errors.New("job repository dependency is required")
	}
	if dependencies.Notification == nil {
		return nil, errors.New("notification dependency is required")
	}

	repo, err := repository.NewRepository(dependencies.DB)
	if err != nil {
		return nil, err
	}

	questionsRepo, err := questionrepository.NewRepository(dependencies.DB)
	if err != nil {
		return nil, err
	}

	userRepo, err := userrepository.NewRepository(dependencies.DB)
	if err != nil {
		return nil, err
	}

	svc, err := service.NewService(
		repo,
		dependencies.JobRepo,
		questionsRepo,
		userRepo,
		dependencies.Notification,
		dependencies.Scheduler,
	)
	if err != nil {
		return nil, err
	}

	ctrl, err := controller.NewController(dependencies.HTTPServer, svc)
	if err != nil {
		return nil, err
	}

	dependencies.Scheduler.RegisterTask(service.TaskNameSendReminder, svc.HandleSendReminderJob)
	dependencies.Scheduler.RegisterTask(service.TaskNamePlanSmartReminders, svc.HandlePlannerJob)

	if err := svc.EnsureNextPlannerJob(ctx); err != nil {
		return nil, err
	}

	return &Module{
		Repository: repo,
		Service:    svc,
		Controller: ctrl,
	}, nil
}
