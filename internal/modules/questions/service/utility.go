package service

import (
	"context"

	"github.com/malcolmkzh/study-notifier/internal/modules/questions/dto"
)

type Utility interface {
	GenerateQuestions(ctx context.Context, noteID uint, userID string) (*dto.GenerateQuestionsResponse, error)
}
