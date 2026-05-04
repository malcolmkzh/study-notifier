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

type sendTelegramQuizPollBody struct {
	ChatID          string   `json:"chat_id"`
	Question        string   `json:"question"`
	Options         []string `json:"options"`
	Type            string   `json:"type"`
	IsAnonymous     bool     `json:"is_anonymous"`
	AllowsRevoting  bool     `json:"allows_revoting"`
	CorrectOptionID int      `json:"correct_option_id"`
	Explanation     string   `json:"explanation,omitempty"`
}

type sendTelegramQuizPollResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Poll struct {
			ID string `json:"id"`
		} `json:"poll"`
	} `json:"result"`
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

func (n *Implementation) SendTelegramQuizPoll(ctx context.Context, request SendTelegramQuizPollRequest) (string, error) {
	if strings.TrimSpace(request.ChatID) == "" {
		return "", errors.New("chat id is required")
	}
	if strings.TrimSpace(request.Question) == "" {
		return "", errors.New("poll question is required")
	}
	if len(request.Options) < 2 {
		return "", errors.New("poll requires at least two options")
	}
	if request.CorrectOptionID < 0 || request.CorrectOptionID >= len(request.Options) {
		return "", errors.New("correct option id is invalid")
	}

	cfg := n.config.Config()
	if strings.TrimSpace(cfg.TelegramBotToken) == "" {
		n.logger.WarnContext(ctx, "telegram bot token is not configured, skipping quiz poll", "chat_id", request.ChatID)
		return "", nil
	}

	baseURL := strings.TrimRight(cfg.TelegramBotBaseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.telegram.org"
	}

	var response sendTelegramQuizPollResponse
	_, err := n.httpClient.DoJSON(ctx, httpclient.Request{
		Method:  http.MethodPost,
		URL:     fmt.Sprintf("%s/bot%s/sendPoll", baseURL, cfg.TelegramBotToken),
		Timeout: 15 * time.Second,
		Body: sendTelegramQuizPollBody{
			ChatID:          request.ChatID,
			Question:        request.Question,
			Options:         request.Options,
			Type:            "quiz",
			IsAnonymous:     false,
			AllowsRevoting:  false,
			CorrectOptionID: request.CorrectOptionID,
			Explanation:     request.Explanation,
		},
	}, &response)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(response.Result.Poll.ID) == "" {
		return "", errors.New("telegram quiz poll response did not contain poll id")
	}

	n.logger.InfoContext(ctx, "telegram quiz poll sent", "chat_id", request.ChatID, "poll_id", response.Result.Poll.ID)
	return response.Result.Poll.ID, nil
}
