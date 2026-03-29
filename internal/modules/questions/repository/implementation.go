package repository

import (
	"context"
	"errors"

	"github.com/malcolmkzh/study-notifier/internal/modules/questions/model"
	"github.com/malcolmkzh/study-notifier/internal/utilities/db"
)

type Implementation struct {
	db db.Utility
}

func NewRepository(dbUtility db.Utility) (*Implementation, error) {
	if dbUtility == nil {
		return nil, errors.New("db utility is required")
	}

	return &Implementation{
		db: dbUtility,
	}, nil
}

func (m *Implementation) CreateMany(ctx context.Context, questions []model.Question) error {
	if len(questions) == 0 {
		return nil
	}

	return m.db.DB().WithContext(ctx).Create(&questions).Error
}
