package notification

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/malcolmkzh/study-notifier/internal/utilities/config"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpclient"
)

type Implementation struct {
	logger     *slog.Logger
	config     config.Utility
	httpClient httpclient.Utility
}

type sendTelegramMessageBody struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func NewNotificationUtility(configUtility config.Utility, httpClientUtility httpclient.Utility) (*Implementation, error) {
	if configUtility == nil {
		return nil, errors.New("config utility is required")
	}
	if httpClientUtility == nil {
		return nil, errors.New("http client utility is required")
	}

	return &Implementation{
		logger:     slog.Default(),
		config:     configUtility,
		httpClient: httpClientUtility,
	}, nil
}

func (n *Implementation) SendTelegramMessage(ctx context.Context, request SendTelegramMessageRequest) error {
	if strings.TrimSpace(request.ChatID) == "" {
		return errors.New("chat id is required")
	}
	if strings.TrimSpace(request.Text) == "" {
		return errors.New("message text is required")
	}

	cfg := n.config.Config()
	if strings.TrimSpace(cfg.TelegramBotToken) == "" {
		n.logger.WarnContext(ctx, "telegram bot token is not configured, skipping send", "chat_id", request.ChatID)
		return nil
	}

	baseURL := strings.TrimRight(cfg.TelegramBotBaseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.telegram.org"
	}

	_, err := n.httpClient.Do(ctx, httpclient.Request{
		Method:  http.MethodPost,
		URL:     fmt.Sprintf("%s/bot%s/sendMessage", baseURL, cfg.TelegramBotToken),
		Timeout: 15 * time.Second,
		Body: sendTelegramMessageBody{
			ChatID: request.ChatID,
			Text:   request.Text,
		},
	})
	if err != nil {
		return err
	}

	n.logger.InfoContext(ctx, "telegram message sent", "chat_id", request.ChatID)
	return nil
}
