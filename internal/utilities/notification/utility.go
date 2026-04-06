package notification

import "context"

type SendQuestionRequest struct {
	TelegramChatID string
	QuestionText   string
	Options        []string
	CorrectOption  string
}

type Utility interface {
	SendQuestion(ctx context.Context, request SendQuestionRequest) error
}
