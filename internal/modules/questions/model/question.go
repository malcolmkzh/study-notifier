package model

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID            uint           `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID        string         `json:"user_id" gorm:"column:user_id;type:varchar(64);not null;index"`
	NoteID        uint           `json:"note_id" gorm:"column:note_id;not null;index"`
	QuestionText  string         `json:"question_text" gorm:"column:question_text;type:text;not null"`
	OptionA       string         `json:"option_a" gorm:"column:option_a;type:text;not null"`
	OptionB       string         `json:"option_b" gorm:"column:option_b;type:text;not null"`
	OptionC       string         `json:"option_c" gorm:"column:option_c;type:text;not null"`
	OptionD       string         `json:"option_d" gorm:"column:option_d;type:text;not null"`
	CorrectOption string         `json:"correct_option" gorm:"column:correct_option;type:char(1);not null"`
	SourceType    string         `json:"source_type" gorm:"column:source_type;type:varchar(50);not null;default:llm"`
	CreatedAt     time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
}

func (Question) TableName() string {
	return "questions"
}
