package reminder

import (
	"context"
	"errors"

	questionrepository "github.com/malcolmkzh/study-notifier/internal/modules/questions/repository"
	"github.com/malcolmkzh/study-notifier/internal/modules/reminder/repository"
	"github.com/malcolmkzh/study-notifier/internal/modules/reminder/service"
	"github.com/malcolmkzh/study-notifier/internal/utilities/db"
	"github.com/malcolmkzh/study-notifier/internal/utilities/notification"
	"github.com/malcolmkzh/study-notifier/internal/utilities/scheduler"
)

type Dependencies struct {
	DB           db.Utility
	Scheduler    scheduler.Utility
	JobRepo      scheduler.JobRepository
	Notification notification.Utility
}

type Module struct {
	Repository repository.Utility
	Service    service.Service
}

func New(ctx context.Context, dependencies Dependencies) (*Module, error) {
	_ = ctx

	if dependencies.DB == nil {
		return nil, errors.New("db dependency is required")
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

	svc, err := service.NewService(
		repo,
		dependencies.JobRepo,
		questionsRepo,
		dependencies.Notification,
		dependencies.Scheduler,
	)
	if err != nil {
		return nil, err
	}

	dependencies.Scheduler.RegisterTask(service.TaskNameSendReminder, svc.HandleScheduledJob)

	return &Module{
		Repository: repo,
		Service:    svc,
	}, nil
}
