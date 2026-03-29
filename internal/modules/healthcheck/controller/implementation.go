package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/malcolmkzh/study-notifier/internal/utilities/httpserver"

	"github.com/gin-gonic/gin"
)

type Implementation struct{}

func NewController(httpServerUtility httpserver.Utility) (*Implementation, error) {
	if httpServerUtility == nil {
		return nil, errors.New("http server utility is required")
	}

	m := &Implementation{}

	err := httpServerUtility.RegisterEndpoint(context.Background(), httpserver.RegisterEndpointRequest{
		Method: http.MethodGet,
		Path:   "/health",
		Fn:     m.HealthCheck,
	})
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Implementation) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "health check endpoint is working",
	})
}
