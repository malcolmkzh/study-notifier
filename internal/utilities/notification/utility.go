package notification

import "context"

type SendTelegramMessageRequest struct {
	ChatID string
	Text   string
}

type SendTelegramQuizPollRequest struct {
	ChatID          string
	Question        string
	Options         []string
	CorrectOptionID int
	Explanation     string
}

type Utility interface {
	SendTelegramMessage(ctx context.Context, request SendTelegramMessageRequest) error
	SendTelegramQuizPoll(ctx context.Context, request SendTelegramQuizPollRequest) (string, error)
}
