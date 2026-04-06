package repository

import (
	"context"
	"errors"

	"github.com/malcolmkzh/study-notifier/internal/modules/telegram/model"
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

func (r *Implementation) Create(ctx context.Context, link *model.TelegramLink) error {
	if link == nil {
		return errors.New("telegram link is required")
	}

	return r.db.DB().WithContext(ctx).Create(link).Error
}

func (r *Implementation) GetByCode(ctx context.Context, code string) (*model.TelegramLink, error) {
	var link model.TelegramLink

	if err := r.db.DB().WithContext(ctx).Where("code = ?", code).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &link, nil
}

func (r *Implementation) DeleteByCode(ctx context.Context, code string) error {
	return r.db.DB().WithContext(ctx).Where("code = ?", code).Delete(&model.TelegramLink{}).Error
}

func (r *Implementation) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.DB().WithContext(ctx).Where("user_id = ?", userID).Delete(&model.TelegramLink{}).Error
}
