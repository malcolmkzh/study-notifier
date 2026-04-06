package model

import "time"

type Account struct {
	ID             string    `json:"id" gorm:"column:id;primaryKey"`
	Email          string    `json:"email" gorm:"column:email;type:varchar(255);not null;uniqueIndex"`
	PasswordHash   string    `json:"-" gorm:"column:password_hash;type:varchar(255);not null"`
	Role           string    `json:"role" gorm:"column:role;type:varchar(50);not null"`
	TelegramChatID *string   `json:"telegram_chat_id,omitempty" gorm:"column:telegram_chat_id;type:varchar(64)"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

func (Account) TableName() string {
	return "accounts"
}
