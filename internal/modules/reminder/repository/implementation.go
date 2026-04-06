package repository

import (
	"context"
	"errors"
	"time"

	"github.com/malcolmkzh/study-notifier/internal/modules/reminder/model"
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

func (r *Implementation) Create(ctx context.Context, reminder *model.Reminder) error {
	if reminder == nil {
		return errors.New("reminder is required")
	}

	return r.db.DB().WithContext(ctx).Create(reminder).Error
}

func (r *Implementation) GetByID(ctx context.Context, id int64) (*model.Reminder, error) {
	var reminder model.Reminder

	if err := r.db.DB().WithContext(ctx).First(&reminder, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &reminder, nil
}

func (r *Implementation) MarkSent(ctx context.Context, id int64) error {
	return r.updateStatus(ctx, id, model.ReminderStatusSent)
}

func (r *Implementation) MarkCancelled(ctx context.Context, id int64) error {
	return r.updateStatus(ctx, id, model.ReminderStatusCancelled)
}

func (r *Implementation) updateStatus(ctx context.Context, id int64, status model.ReminderStatus) error {
	return r.db.DB().WithContext(ctx).
		Model(&model.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now().UTC(),
		}).Error
}
