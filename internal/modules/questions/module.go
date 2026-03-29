package questions

import (
	"context"
	"errors"

	notesrepository "github.com/malcolmkzh/study-notifier/internal/modules/notes/repository"
	"github.com/malcolmkzh/study-notifier/internal/modules/questions/controller"
	"github.com/malcolmkzh/study-notifier/internal/modules/questions/repository"
	"github.com/malcolmkzh/study-notifier/internal/modules/questions/service"
	"github.com/malcolmkzh/study-notifier/internal/utilities/db"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpserver"
	"github.com/malcolmkzh/study-notifier/internal/utilities/llm"
)

type Dependencies struct {
	DB         db.Utility
	LLM        llm.Utility
	HTTPServer httpserver.Utility
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
	if dependencies.LLM == nil {
		return nil, errors.New("llm dependency is required")
	}

	repo, err := repository.NewRepository(dependencies.DB)
	if err != nil {
		return nil, err
	}

	notesRepo, err := notesrepository.NewRepository(dependencies.DB)
	if err != nil {
		return nil, err
	}

	svc, err := service.NewService(repo, notesRepo, dependencies.LLM)
	if err != nil {
		return nil, err
	}

	ctrl, err := controller.NewController(dependencies.HTTPServer, svc)
	if err != nil {
		return nil, err
	}

	return &Module{
		Repository: repo,
		Service:    svc,
		Controller: ctrl,
	}, nil
}
