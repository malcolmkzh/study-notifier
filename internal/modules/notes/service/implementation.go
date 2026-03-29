package service

import (
	"context"
	"errors"

	"github.com/malcolmkzh/study-notifier/internal/modules/notes/dto"
	"github.com/malcolmkzh/study-notifier/internal/modules/notes/model"
	"github.com/malcolmkzh/study-notifier/internal/modules/notes/repository"
)

type Implementation struct {
	repository repository.Utility
}

func NewService(repo repository.Utility) (*Implementation, error) {
	if repo == nil {
		return nil, errors.New("notes repository is required")
	}

	return &Implementation{
		repository: repo,
	}, nil
}

func (m *Implementation) Create(ctx context.Context, userID string, request dto.CreateNoteRequest) (*model.Note, error) {
	return m.repository.Create(ctx, model.Note{
		UserID:  userID,
		Title:   request.Title,
		Content: request.Content,
		Topic:   request.Topic,
		Tags:    request.Tags,
	})
}

func (m *Implementation) List(ctx context.Context, userID string) ([]model.Note, error) {
	return m.repository.ListByUserID(ctx, userID)
}

func (m *Implementation) GetByID(ctx context.Context, id uint, userID string) (*model.Note, error) {
	return m.repository.GetByID(ctx, id, userID)
}

func (m *Implementation) Update(ctx context.Context, id uint, userID string, request dto.UpdateNoteRequest) (*model.Note, error) {
	note, err := m.repository.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, nil
	}

	note.Title = request.Title
	note.Content = request.Content
	note.Topic = request.Topic
	note.Tags = request.Tags

	return m.repository.Update(ctx, *note)
}

func (m *Implementation) Delete(ctx context.Context, id uint, userID string) error {
	return m.repository.Delete(ctx, id, userID)
}
