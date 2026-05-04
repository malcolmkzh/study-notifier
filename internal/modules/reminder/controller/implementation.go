package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	reminderdto "github.com/malcolmkzh/study-notifier/internal/modules/reminder/dto"
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
		{Method: http.MethodPost, Path: "/reminders/plan-today", Fn: controller.PlanToday, RequireAuth: true},
		{Method: http.MethodGet, Path: "/reminder-settings", Fn: controller.GetReminderSetting, RequireAuth: true},
		{Method: http.MethodPut, Path: "/reminder-settings", Fn: controller.UpdateReminderSetting, RequireAuth: true},
	}

	for _, route := range routes {
		if err := httpServerUtility.RegisterEndpoint(context.Background(), route); err != nil {
			return nil, err
		}
	}

	return controller, nil
}

func (c *Implementation) GetReminderSetting(ctx *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(ctx)
	if !ok {
		return
	}

	response, err := c.service.GetReminderSetting(ctx.Request.Context(), userID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Implementation) UpdateReminderSetting(ctx *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(ctx)
	if !ok {
		return
	}

	var request reminderdto.UpdateReminderSettingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	response, err := c.service.UpdateReminderSetting(ctx.Request.Context(), service.UpdateReminderSettingRequest{
		UserID:  userID,
		Enabled: request.Enabled,
	})
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// Test function to run planner immediately without waiting for the cron job. Not exposed in production.
func (c *Implementation) PlanToday(ctx *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(ctx)
	if !ok {
		return
	}

	err := c.service.PlanSmartReminders(ctx.Request.Context(), &userID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "smart reminders planned successfully",
	})
}
