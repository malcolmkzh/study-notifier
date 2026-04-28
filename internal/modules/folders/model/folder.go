package model

import (
	"time"

	"gorm.io/gorm"
)

type Folder struct {
	ID        uint           `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID    string         `json:"user_id" gorm:"column:user_id;index;not null"`
	Name      string         `json:"name" gorm:"column:name;not null"`
	ParentID  *uint          `json:"parent_id" gorm:"column:parent_id;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Folder) TableName() string {
	return "note_folders"
}
