package repository

import (
	"context"

	"github.com/malcolmkzh/study-notifier/internal/modules/notes/model"
)

type Utility interface {
	Create(ctx context.Context, note model.Note) (*model.Note, error)
	ListByUserID(ctx context.Context, userID string) ([]model.Note, error)
	GetByID(ctx context.Context, id uint, userID string) (*model.Note, error)
	Update(ctx context.Context, note model.Note) (*model.Note, error)
	Delete(ctx context.Context, id uint, userID string) error
}
