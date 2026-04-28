package repository

import (
	"context"
	"errors"

	"github.com/malcolmkzh/study-notifier/internal/modules/folders/model"
	notemodel "github.com/malcolmkzh/study-notifier/internal/modules/notes/model"
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

func (m *Implementation) Create(ctx context.Context, folder model.Folder) (*model.Folder, error) {
	if err := m.db.DB().WithContext(ctx).Create(&folder).Error; err != nil {
		return nil, err
	}

	return &folder, nil
}

func (m *Implementation) ListByUserID(ctx context.Context, userID string) ([]model.Folder, error) {
	var folders []model.Folder
	if err := m.db.DB().WithContext(ctx).Where("user_id = ?", userID).Order("name ASC").Find(&folders).Error; err != nil {
		return nil, err
	}

	return folders, nil
}

func (m *Implementation) GetByID(ctx context.Context, id uint, userID string) (*model.Folder, error) {
	var folder model.Folder
	if err := m.db.DB().WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&folder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &folder, nil
}

func (m *Implementation) Update(ctx context.Context, folder model.Folder) (*model.Folder, error) {
	if err := m.db.DB().WithContext(ctx).Save(&folder).Error; err != nil {
		return nil, err
	}

	return &folder, nil
}

func (m *Implementation) Delete(ctx context.Context, id uint, userID string) error {
	return m.db.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&notemodel.Note{}).Where("folder_id = ? AND user_id = ?", id, userID).Update("folder_id", gorm.Expr("NULL"))
		if result.Error != nil {
			return result.Error
		}

		result = tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Folder{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}
