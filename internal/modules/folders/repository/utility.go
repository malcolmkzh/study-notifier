package repository

import (
	"context"

	"github.com/malcolmkzh/study-notifier/internal/modules/folders/model"
)

type Utility interface {
	Create(ctx context.Context, folder model.Folder) (*model.Folder, error)
	ListByUserID(ctx context.Context, userID string) ([]model.Folder, error)
	GetByID(ctx context.Context, id uint, userID string) (*model.Folder, error)
	Update(ctx context.Context, folder model.Folder) (*model.Folder, error)
	Delete(ctx context.Context, id uint, userID string) error
}
