package repository

import (
	"context"

	"github.com/malcolmkzh/study-notifier/internal/modules/questions/model"
)

type Utility interface {
	CreateMany(ctx context.Context, questions []model.Question) error
	GetRandomQuestionByUserID(ctx context.Context, userID string) (*model.Question, error)
}
