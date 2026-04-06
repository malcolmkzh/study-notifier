package notification

import "context"

type SendTelegramMessageRequest struct {
	ChatID string
	Text   string
}

type Utility interface {
	SendTelegramMessage(ctx context.Context, request SendTelegramMessageRequest) error
}
