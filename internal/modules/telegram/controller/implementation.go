package controller

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/malcolmkzh/study-notifier/internal/modules/telegram/service"
	"github.com/malcolmkzh/study-notifier/internal/utilities/config"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpserver"
)

type Implementation struct {
	httpServerUtility httpserver.Utility
	service           service.Utility
	configUtility     config.Utility
}

type telegramUpdate struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func NewController(httpServerUtility httpserver.Utility, service service.Utility, configUtility config.Utility) (*Implementation, error) {
	if httpServerUtility == nil {
		return nil, errors.New("http server utility is required")
	}
	if service == nil {
		return nil, errors.New("telegram service is required")
	}
	if configUtility == nil {
		return nil, errors.New("config utility is required")
	}

	controller := &Implementation{
		httpServerUtility: httpServerUtility,
		service:           service,
		configUtility:     configUtility,
	}

	routes := []httpserver.RegisterEndpointRequest{
		{Method: http.MethodPost, Path: "/telegram/link", Fn: controller.CreateLink, RequireAuth: true},
		{Method: http.MethodPost, Path: "/telegram/webhook", Fn: controller.Webhook},
	}

	for _, route := range routes {
		if err := httpServerUtility.RegisterEndpoint(context.Background(), route); err != nil {
			return nil, err
		}
	}

	return controller, nil
}

func (c *Implementation) CreateLink(ctx *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(ctx)
	if !ok {
		return
	}

	response, err := c.service.CreateLink(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Implementation) Webhook(ctx *gin.Context) {
	expectedSecret := strings.TrimSpace(c.configUtility.Config().TelegramWebhookSecret)
	if expectedSecret != "" && ctx.GetHeader("X-Telegram-Bot-Api-Secret-Token") != expectedSecret {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var update telegramUpdate
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx.Status(http.StatusOK)

	if update.Message.Text == "" {
		return
	}

	go func(chatID int64, text string) {
		_ = c.service.HandleMessage(context.Background(), chatID, text)
	}(update.Message.Chat.ID, update.Message.Text)
}
