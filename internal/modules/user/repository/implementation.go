package repository

import (
	"context"
	"errors"
	"time"

	"github.com/malcolmkzh/study-notifier/internal/modules/user/model"
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

func (r *Implementation) GetByID(ctx context.Context, userID string) (*model.Account, error) {
	var account model.Account

	if err := r.db.DB().WithContext(ctx).Where("id = ?", userID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}

func (r *Implementation) UpdateTelegramChatID(ctx context.Context, userID string, chatID string) error {
	return r.db.DB().WithContext(ctx).
		Model(&model.Account{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"telegram_chat_id": chatID,
			"updated_at":       time.Now().UTC(),
		}).Error
}
