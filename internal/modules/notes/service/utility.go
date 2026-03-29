package service

import (
	"context"

	"github.com/malcolmkzh/study-notifier/internal/modules/notes/dto"
	"github.com/malcolmkzh/study-notifier/internal/modules/notes/model"
)

type Utility interface {
	Create(ctx context.Context, userID string, request dto.CreateNoteRequest) (*model.Note, error)
	List(ctx context.Context, userID string) ([]model.Note, error)
	GetByID(ctx context.Context, id uint, userID string) (*model.Note, error)
	Update(ctx context.Context, id uint, userID string, request dto.UpdateNoteRequest) (*model.Note, error)
	Delete(ctx context.Context, id uint, userID string) error
}
