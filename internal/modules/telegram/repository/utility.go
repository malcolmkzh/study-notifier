package repository

import (
	"context"

	"github.com/malcolmkzh/study-notifier/internal/modules/telegram/model"
)

type Utility interface {
	Create(ctx context.Context, link *model.TelegramLink) error
	GetByCode(ctx context.Context, code string) (*model.TelegramLink, error)
	DeleteByCode(ctx context.Context, code string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
