package repository

import (
	"context"
	"errors"

	"github.com/malcolmkzh/study-notifier/internal/modules/questions/model"
	"github.com/malcolmkzh/study-notifier/internal/utilities/db"
	"gorm.io/gorm"
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

func (m *Implementation) GetRandomQuestionByUserID(ctx context.Context, userID string) (*model.Question, error) {
	var question model.Question

	if err := m.db.DB().WithContext(ctx).
		Where("user_id = ?", userID).
		Order("RAND()").
		First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &question, nil
}
