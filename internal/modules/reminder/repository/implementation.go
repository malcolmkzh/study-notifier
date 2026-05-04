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

//Get pending reminders with scheduled_at between the specified start and end time. Used to prevent over-scheduling and for testing daily planning
func (r *Implementation) ListPendingByUserIDBetween(ctx context.Context, userID string, start time.Time, end time.Time) ([]model.Reminder, error) {
	var reminders []model.Reminder

	if err := r.db.DB().WithContext(ctx).
		Where("user_id = ? AND status = ? AND scheduled_at >= ? AND scheduled_at < ?", userID, model.ReminderStatusPending, start, end).
		Order("scheduled_at ASC").
		Find(&reminders).Error; err != nil {
		return nil, err
	}

	return reminders, nil
}

//Get reminders with scheduled_at after the specified time, regardless of the end time. Used for fetching upcoming reminders for deleting
func (r *Implementation) ListPendingByUserIDAfter(ctx context.Context, userID string, after time.Time) ([]model.Reminder, error) {
	var reminders []model.Reminder

	if err := r.db.DB().WithContext(ctx).
		Where("user_id = ? AND status = ? AND scheduled_at >= ?", userID, model.ReminderStatusPending, after).
		Order("scheduled_at ASC").
		Find(&reminders).Error; err != nil {
		return nil, err
	}

	return reminders, nil
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

func (r *Implementation) GetSettingByUserID(ctx context.Context, userID string) (*model.ReminderSetting, error) {
	var setting model.ReminderSetting

	if err := r.db.DB().WithContext(ctx).Where("user_id = ?", userID).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &setting, nil
}

func (r *Implementation) UpsertSetting(ctx context.Context, setting *model.ReminderSetting) error {
	if setting == nil {
		return errors.New("reminder setting is required")
	}

	now := time.Now().UTC()
	return r.db.DB().WithContext(ctx).
		Model(&model.ReminderSetting{}).
		Where("user_id = ?", setting.UserID).
		Assign(map[string]interface{}{
			"enabled":    setting.Enabled,
			"timezone":   setting.Timezone,
			"updated_at": now,
		}).
		FirstOrCreate(&model.ReminderSetting{
			UserID:    setting.UserID,
			Enabled:   setting.Enabled,
			Timezone:  setting.Timezone,
			CreatedAt: now,
			UpdatedAt: now,
		}).Error
}

func (r *Implementation) ListEnabledSettings(ctx context.Context) ([]model.ReminderSetting, error) {
	var settings []model.ReminderSetting

	if err := r.db.DB().WithContext(ctx).
		Where("enabled = ?", true).
		Find(&settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}

func (r *Implementation) CreateQuestionAttempt(ctx context.Context, attempt *model.QuestionAttempt) error {
	if attempt == nil {
		return errors.New("question attempt is required")
	}

	return r.db.DB().WithContext(ctx).Create(attempt).Error
}

func (r *Implementation) GetQuestionAttemptByTelegramPollID(ctx context.Context, pollID string) (*model.QuestionAttempt, error) {
	var attempt model.QuestionAttempt

	if err := r.db.DB().WithContext(ctx).Where("telegram_poll_id = ?", pollID).First(&attempt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &attempt, nil
}

func (r *Implementation) ListAnsweredAttemptsByUserID(ctx context.Context, userID string) ([]model.QuestionAttempt, error) {
	var attempts []model.QuestionAttempt

	if err := r.db.DB().WithContext(ctx).
		Where("user_id = ? AND answered_at IS NOT NULL", userID).
		Order("answered_at DESC").
		Limit(100).
		Find(&attempts).Error; err != nil {
		return nil, err
	}

	return attempts, nil
}

func (r *Implementation) UpdateQuestionAttemptTelegramPollID(ctx context.Context, id int64, pollID string) error {
	return r.db.DB().WithContext(ctx).
		Model(&model.QuestionAttempt{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"telegram_poll_id": pollID,
			"updated_at":       time.Now().UTC(),
		}).Error
}

func (r *Implementation) MarkQuestionAttemptPollAnswered(ctx context.Context, id int64, answeredAt time.Time, isCorrect bool) error {
	return r.db.DB().WithContext(ctx).
		Model(&model.QuestionAttempt{}).
		Where("id = ? AND answered_at IS NULL", id).
		Updates(map[string]interface{}{
			"answered_at": answeredAt.UTC(),
			"is_correct":  isCorrect,
			"updated_at":  time.Now().UTC(),
		}).Error
}
