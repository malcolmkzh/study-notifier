package model

import "time"

type Note struct {
	ID        uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID    string    `json:"user_id" gorm:"index;not null"`
	FolderID  *uint     `json:"folder_id" gorm:"column:folder_id;index"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Topic     string    `json:"topic"`
	Tags      []string  `json:"tags" gorm:"serializer:json"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Note) TableName() string {
	return "notes"
}
