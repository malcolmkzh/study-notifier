package telegram

import (
	"context"
	"errors"

	"github.com/malcolmkzh/study-notifier/internal/modules/telegram/controller"
	"github.com/malcolmkzh/study-notifier/internal/modules/telegram/repository"
	"github.com/malcolmkzh/study-notifier/internal/modules/telegram/service"
	userrepository "github.com/malcolmkzh/study-notifier/internal/modules/user/repository"
	"github.com/malcolmkzh/study-notifier/internal/utilities/config"
	"github.com/malcolmkzh/study-notifier/internal/utilities/db"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpserver"
	"github.com/malcolmkzh/study-notifier/internal/utilities/notification"
)

type Dependencies struct {
	DB           db.Utility
	HTTPServer   httpserver.Utility
	Notification notification.Utility
	Config       config.Utility
}

type Module struct {
	Repository repository.Utility
	Service    service.Utility
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
	if dependencies.Notification == nil {
		return nil, errors.New("notification dependency is required")
	}
	if dependencies.Config == nil {
		return nil, errors.New("config dependency is required")
	}

	repo, err := repository.NewRepository(dependencies.DB)
	if err != nil {
		return nil, err
	}

	userRepo, err := userrepository.NewRepository(dependencies.DB)
	if err != nil {
		return nil, err
	}

	svc, err := service.NewService(repo, userRepo, dependencies.Notification)
	if err != nil {
		return nil, err
	}

	ctrl, err := controller.NewController(dependencies.HTTPServer, svc, dependencies.Config)
	if err != nil {
		return nil, err
	}

	return &Module{
		Repository: repo,
		Service:    svc,
		Controller: ctrl,
	}, nil
}
