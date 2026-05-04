package model

import "time"

type QuestionAttempt struct {
	ID              int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID          string     `json:"user_id" gorm:"column:user_id;type:varchar(64);not null;index"`
	ReminderID      int64      `json:"reminder_id" gorm:"column:reminder_id;not null;index"`
	QuestionID      uint       `json:"question_id" gorm:"column:question_id;not null;index"`
	TelegramPollID  *string    `json:"telegram_poll_id" gorm:"column:telegram_poll_id;type:varchar(128);uniqueIndex"`
	CorrectOptionID int        `json:"correct_option_id" gorm:"column:correct_option_id;not null"`
	SentAt          time.Time  `json:"sent_at" gorm:"column:sent_at;not null"`
	SelectedOption  *string    `json:"selected_option" gorm:"column:selected_option;type:char(1)"`
	AnsweredAt      *time.Time `json:"answered_at" gorm:"column:answered_at"`
	IsCorrect       *bool      `json:"is_correct" gorm:"column:is_correct"`
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"column:updated_at;not null"`
}

func (QuestionAttempt) TableName() string {
	return "question_attempts"
}
