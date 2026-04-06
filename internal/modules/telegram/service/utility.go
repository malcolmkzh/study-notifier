package service

import (
	"context"

	"github.com/malcolmkzh/study-notifier/internal/modules/telegram/dto"
)

type Utility interface {
	CreateLink(ctx context.Context, userID string) (*dto.CreateLinkResponse, error)
	HandleMessage(ctx context.Context, chatID int64, text string) error
}
