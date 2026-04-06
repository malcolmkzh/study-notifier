package model

import "time"

type TelegramLink struct {
	ID        uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID    string    `json:"user_id" gorm:"column:user_id;type:varchar(64);not null;index"`
	Code      string    `json:"code" gorm:"column:code;type:varchar(32);not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"column:expires_at;not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null"`
}

func (TelegramLink) TableName() string {
	return "telegram_links"
}
