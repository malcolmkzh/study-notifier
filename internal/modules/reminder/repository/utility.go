package repository

import (
	"context"

	"github.com/malcolmkzh/study-notifier/internal/modules/reminder/model"
)

type Utility interface {
	Create(ctx context.Context, reminder *model.Reminder) error
	GetByID(ctx context.Context, id int64) (*model.Reminder, error)
	MarkSent(ctx context.Context, id int64) error
	MarkCancelled(ctx context.Context, id int64) error
}
