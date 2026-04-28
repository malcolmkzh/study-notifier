package service

import (
	"context"
	"errors"

	folderrepository "github.com/malcolmkzh/study-notifier/internal/modules/folders/repository"
	"github.com/malcolmkzh/study-notifier/internal/modules/notes/dto"
	"github.com/malcolmkzh/study-notifier/internal/modules/notes/model"
	"github.com/malcolmkzh/study-notifier/internal/modules/notes/repository"
	"github.com/malcolmkzh/study-notifier/internal/utilities/errorutil"
	"gorm.io/gorm"
)

type Implementation struct {
	repository       repository.Utility
	folderRepository folderrepository.Utility
}

func NewService(repo repository.Utility, folderRepo folderrepository.Utility) (*Implementation, error) {
	if repo == nil {
		return nil, errors.New("notes repository is required")
	}
	if folderRepo == nil {
		return nil, errors.New("folder repository is required")
	}

	return &Implementation{
		repository:       repo,
		folderRepository: folderRepo,
	}, nil
}

func (m *Implementation) Create(ctx context.Context, userID string, request dto.CreateNoteRequest) (*model.Note, error) {
	if err := m.validateFolder(ctx, userID, request.FolderID); err != nil {
		return nil, err
	}

	note, err := m.repository.Create(ctx, model.Note{
		UserID:   userID,
		FolderID: request.FolderID,
		Title:    request.Title,
		Content:  request.Content,
		Topic:    request.Topic,
		Tags:     request.Tags,
	})
	if err != nil {
		return nil, err
	}

	return note, nil
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
		return nil, errorutil.NewWithMessage(errorutil.CodeNotFound, "note not found")
	}

	note.Title = request.Title
	note.Content = request.Content
	note.Topic = request.Topic
	note.Tags = request.Tags
	note.FolderID = request.FolderID

	if err := m.validateFolder(ctx, userID, request.FolderID); err != nil {
		return nil, err
	}

	updatedNote, err := m.repository.Update(ctx, *note)
	if err != nil {
		return nil, err
	}

	return updatedNote, nil
}

func (m *Implementation) Delete(ctx context.Context, id uint, userID string) error {
	if err := m.repository.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorutil.NewWithMessage(errorutil.CodeNotFound, "note not found")
		}
		return err
	}

	return nil
}

func (m *Implementation) validateFolder(ctx context.Context, userID string, folderID *uint) error {
	if folderID == nil {
		return nil
	}

	folder, err := m.folderRepository.GetByID(ctx, *folderID, userID)
	if err != nil {
		return err
	}
	if folder == nil {
		return errorutil.NewWithMessage(errorutil.CodeBadRequest, "folder not found")
	}

	return nil
}
