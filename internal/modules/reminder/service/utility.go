package service

import (
	"context"
	"encoding/json"
	"time"
)

type Service interface {
	CreateReminder(ctx context.Context, req CreateReminderRequest) error
	CancelReminder(ctx context.Context, reminderID int64) error
	TriggerReminder(ctx context.Context, reminderID int64) error
	HandleScheduledJob(ctx context.Context, metadata json.RawMessage) error
}

type CreateReminderRequest struct {
	UserID      string
	ScheduledAt time.Time
}
