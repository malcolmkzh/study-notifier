package service

import (
	"context"
	"encoding/json"
	"time"

	reminderdto "github.com/malcolmkzh/study-notifier/internal/modules/reminder/dto"
)

type Service interface {
	GetReminderSetting(ctx context.Context, userID string) (*reminderdto.ReminderSettingResponse, error)
	UpdateReminderSetting(ctx context.Context, req UpdateReminderSettingRequest) (*reminderdto.ReminderSettingResponse, error)
	CreateReminder(ctx context.Context, req CreateReminderRequest) error
	PlanSmartReminders(ctx context.Context, userID *string) error
	TriggerReminder(ctx context.Context, reminderID int64) error
	HandleSendReminderJob(ctx context.Context, metadata json.RawMessage) error
	HandlePlannerJob(ctx context.Context, metadata json.RawMessage) error
	HandlePollAnswer(ctx context.Context, pollID string, optionIDs []int) error
}

type CreateReminderRequest struct {
	UserID      string
	ScheduledAt time.Time
}

type UpdateReminderSettingRequest struct {
	UserID  string
	Enabled bool
}
