package repository

import (
	"context"

	"github.com/malcolmkzh/study-notifier/internal/modules/user/model"
)

type Utility interface {
	GetByID(ctx context.Context, userID string) (*model.Account, error)
	UpdateTelegramChatID(ctx context.Context, userID string, chatID string) error
}
