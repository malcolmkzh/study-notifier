package repository

import (
	"context"
	"time"

	"github.com/malcolmkzh/study-notifier/internal/modules/reminder/model"
)

type Utility interface {
	Create(ctx context.Context, reminder *model.Reminder) error
	GetByID(ctx context.Context, id int64) (*model.Reminder, error)
	ListPendingByUserIDBetween(ctx context.Context, userID string, start time.Time, end time.Time) ([]model.Reminder, error)
	ListPendingByUserIDAfter(ctx context.Context, userID string, after time.Time) ([]model.Reminder, error)
	MarkSent(ctx context.Context, id int64) error
	MarkCancelled(ctx context.Context, id int64) error
	GetSettingByUserID(ctx context.Context, userID string) (*model.ReminderSetting, error)
	UpsertSetting(ctx context.Context, setting *model.ReminderSetting) error
	ListEnabledSettings(ctx context.Context) ([]model.ReminderSetting, error)
	CreateQuestionAttempt(ctx context.Context, attempt *model.QuestionAttempt) error
	GetQuestionAttemptByTelegramPollID(ctx context.Context, pollID string) (*model.QuestionAttempt, error)
	ListAnsweredAttemptsByUserID(ctx context.Context, userID string) ([]model.QuestionAttempt, error)
	UpdateQuestionAttemptTelegramPollID(ctx context.Context, id int64, pollID string) error
	MarkQuestionAttemptPollAnswered(ctx context.Context, id int64, answeredAt time.Time, isCorrect bool) error
}
