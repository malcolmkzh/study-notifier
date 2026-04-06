package model

import "time"

type ReminderStatus string

const (
	ReminderStatusPending   ReminderStatus = "pending"
	ReminderStatusSent      ReminderStatus = "sent"
	ReminderStatusCancelled ReminderStatus = "cancelled"
)

type Reminder struct {
	ID             int64          `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID         string         `json:"user_id" gorm:"column:user_id;type:varchar(64);not null;index"`
	TelegramChatID string         `json:"telegram_chat_id" gorm:"column:telegram_chat_id;type:varchar(64);not null"`
	ScheduledAt    time.Time      `json:"scheduled_at" gorm:"column:scheduled_at;not null;index"`
	Status         ReminderStatus `json:"status" gorm:"column:status;type:varchar(50);not null;index"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
}

func (Reminder) TableName() string {
	return "reminders"
}
