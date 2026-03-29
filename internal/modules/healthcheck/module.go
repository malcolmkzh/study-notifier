package healthcheck

import (
	"context"
	"errors"

	"github.com/malcolmkzh/study-notifier/internal/modules/healthcheck/controller"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpserver"
)

type Dependencies struct {
	HTTPServer httpserver.Utility
}

type Module struct {
	Controller *controller.Implementation
}

func New(ctx context.Context, dependencies Dependencies) (*Module, error) {
	_ = ctx

	if dependencies.HTTPServer == nil {
		return nil, errors.New("http server dependency is required")
	}

	ctrl, err := controller.NewController(dependencies.HTTPServer)
	if err != nil {
		return nil, err
	}

	return &Module{
		Controller: ctrl,
	}, nil
}
