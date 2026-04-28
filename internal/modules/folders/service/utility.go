package service

import (
	"context"

	"github.com/malcolmkzh/study-notifier/internal/modules/folders/dto"
	"github.com/malcolmkzh/study-notifier/internal/modules/folders/model"
)

type Utility interface {
	Create(ctx context.Context, userID string, request dto.CreateFolderRequest) (*model.Folder, error)
	List(ctx context.Context, userID string) ([]model.Folder, error)
	Update(ctx context.Context, id uint, userID string, request dto.UpdateFolderRequest) (*model.Folder, error)
	Delete(ctx context.Context, id uint, userID string) error
}
