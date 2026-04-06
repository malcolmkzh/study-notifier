package controller

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/malcolmkzh/study-notifier/internal/modules/reminder/service"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpserver"
)

type Implementation struct {
	httpServerUtility httpserver.Utility
	service           service.Service
}

func NewController(httpServerUtility httpserver.Utility, service service.Service) (*Implementation, error) {
	if httpServerUtility == nil {
		return nil, errors.New("http server utility is required")
	}
	if service == nil {
		return nil, errors.New("reminder service is required")
	}

	controller := &Implementation{
		httpServerUtility: httpServerUtility,
		service:           service,
	}

	routes := []httpserver.RegisterEndpointRequest{
		{Method: http.MethodPost, Path: "/reminders/test", Fn: controller.CreateTestReminder, RequireAuth: true},
	}

	for _, route := range routes {
		if err := httpServerUtility.RegisterEndpoint(context.Background(), route); err != nil {
			return nil, err
		}
	}

	return controller, nil
}

func (c *Implementation) CreateTestReminder(ctx *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(ctx)
	if !ok {
		return
	}

	scheduledAt := time.Now().UTC().Add(2 * time.Minute)
	err := c.service.CreateReminder(ctx.Request.Context(), service.CreateReminderRequest{
		UserID:      userID,
		ScheduledAt: scheduledAt,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":      "test reminder created successfully",
		"scheduled_at": scheduledAt,
	})
}
