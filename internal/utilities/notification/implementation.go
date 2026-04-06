package notification

import (
	"context"
	"log/slog"
)

type Implementation struct {
	logger *slog.Logger
}

func NewNotificationUtility() *Implementation {
	return &Implementation{
		logger: slog.Default(),
	}
}

func (n *Implementation) SendQuestion(ctx context.Context, request SendQuestionRequest) error {
	n.logger.InfoContext(ctx,
		"notification send question stub",
		"telegram_chat_id", request.TelegramChatID,
		"question_text", request.QuestionText,
	)

	return nil
}
