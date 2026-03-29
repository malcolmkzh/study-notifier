package db

import "gorm.io/gorm"

type Utility interface {
	DB() *gorm.DB
}
