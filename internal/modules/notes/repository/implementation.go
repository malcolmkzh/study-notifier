package repository

import (
	"context"
	"errors"

	"github.com/malcolmkzh/study-notifier/internal/modules/notes/model"
	"github.com/malcolmkzh/study-notifier/internal/utilities/db"
	"gorm.io/gorm"
)

type Implementation struct {
	db db.Utility
}

func NewRepository(dbUtility db.Utility) (*Implementation, error) {
	if dbUtility == nil {
		return nil, errors.New("db utility is required")
	}

	return &Implementation{
		db: dbUtility,
	}, nil
}

func (m *Implementation) Create(ctx context.Context, note model.Note) (*model.Note, error) {
	if err := m.db.DB().WithContext(ctx).Create(&note).Error; err != nil {
		return nil, err
	}

	return &note, nil
}

func (m *Implementation) ListByUserID(ctx context.Context, userID string) ([]model.Note, error) {
	var notes []model.Note
	if err := m.db.DB().WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func (m *Implementation) GetByID(ctx context.Context, id uint, userID string) (*model.Note, error) {
	var note model.Note
	if err := m.db.DB().WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &note, nil
}

func (m *Implementation) Update(ctx context.Context, note model.Note) (*model.Note, error) {
	if err := m.db.DB().WithContext(ctx).Save(&note).Error; err != nil {
		return nil, err
	}

	return &note, nil
}

func (m *Implementation) Delete(ctx context.Context, id uint, userID string) error {
	result := m.db.DB().WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Note{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
