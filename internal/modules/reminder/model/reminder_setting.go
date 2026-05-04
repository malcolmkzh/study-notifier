package model

import "time"

type ReminderSetting struct {
	UserID    string    `json:"user_id" gorm:"column:user_id;type:varchar(64);primaryKey"`
	Enabled   bool      `json:"enabled" gorm:"column:enabled;not null;default:false"`
	Timezone  string    `json:"timezone" gorm:"column:timezone;type:varchar(100);not null;default:Asia/Singapore"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

func (ReminderSetting) TableName() string {
	return "reminder_settings"
}
