package service

import (
	"context"
	"errors"

	"github.com/malcolmkzh/study-notifier/internal/modules/folders/dto"
	"github.com/malcolmkzh/study-notifier/internal/modules/folders/model"
	"github.com/malcolmkzh/study-notifier/internal/modules/folders/repository"
	"github.com/malcolmkzh/study-notifier/internal/utilities/errorutil"
	"gorm.io/gorm"
)

type Implementation struct {
	repository repository.Utility
}

func NewService(repo repository.Utility) (*Implementation, error) {
	if repo == nil {
		return nil, errors.New("folders repository is required")
	}

	return &Implementation{
		repository: repo,
	}, nil
}

func (m *Implementation) Create(ctx context.Context, userID string, request dto.CreateFolderRequest) (*model.Folder, error) {
	if request.ParentID != nil {
		parent, err := m.repository.GetByID(ctx, *request.ParentID, userID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, errorutil.NewWithMessage(errorutil.CodeBadRequest, "parent folder not found")
		}
	}

	return m.repository.Create(ctx, model.Folder{
		UserID:   userID,
		Name:     request.Name,
		ParentID: request.ParentID,
	})
}

func (m *Implementation) List(ctx context.Context, userID string) ([]model.Folder, error) {
	return m.repository.ListByUserID(ctx, userID)
}

func (m *Implementation) Update(ctx context.Context, id uint, userID string, request dto.UpdateFolderRequest) (*model.Folder, error) {
	folder, err := m.repository.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, errorutil.NewWithMessage(errorutil.CodeNotFound, "folder not found")
	}

	if request.ParentID != nil {
		if *request.ParentID == id {
			return nil, errorutil.NewWithMessage(errorutil.CodeBadRequest, "folder cannot parent itself")
		}

		parent, err := m.repository.GetByID(ctx, *request.ParentID, userID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, errorutil.NewWithMessage(errorutil.CodeBadRequest, "parent folder not found")
		}
	}

	folder.Name = request.Name
	folder.ParentID = request.ParentID

	return m.repository.Update(ctx, *folder)
}

func (m *Implementation) Delete(ctx context.Context, id uint, userID string) error {
	if err := m.repository.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorutil.NewWithMessage(errorutil.CodeNotFound, "folder not found")
		}
		return err
	}

	return nil
}
