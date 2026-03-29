package service

import (
	"context"
	"errors"
	"strings"

	notesrepository "github.com/malcolmkzh/study-notifier/internal/modules/notes/repository"
	"github.com/malcolmkzh/study-notifier/internal/modules/questions/dto"
	"github.com/malcolmkzh/study-notifier/internal/modules/questions/model"
	"github.com/malcolmkzh/study-notifier/internal/modules/questions/repository"
	"github.com/malcolmkzh/study-notifier/internal/utilities/llm"
)

type Implementation struct {
	repository      repository.Utility
	llmUtility      llm.Utility
	notesRepository notesrepository.Utility
}

func NewService(repo repository.Utility, notesRepo notesrepository.Utility, llmUtility llm.Utility) (*Implementation, error) {
	if repo == nil {
		return nil, errors.New("questions repository is required")
	}
	if notesRepo == nil {
		return nil, errors.New("notes repository is required")
	}
	if llmUtility == nil {
		return nil, errors.New("llm utility is required")
	}

	return &Implementation{
		repository:      repo,
		llmUtility:      llmUtility,
		notesRepository: notesRepo,
	}, nil
}

func (m *Implementation) GenerateQuestions(ctx context.Context, noteID uint, userID string) (*dto.GenerateQuestionsResponse, error) {
	note, err := m.notesRepository.GetByID(ctx, noteID, userID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, nil
	}

	questions, err := m.llmUtility.GenerateQuestions(ctx, llm.GenerateQuestionsRequest{
		Title:   note.Title,
		Content: note.Content,
		Topic:   note.Topic,
		Tags:    note.Tags,
	})
	if err != nil {
		return nil, err
	}

	questionModels := make([]model.Question, 0, len(questions.Questions))
	for _, question := range questions.Questions {
		questionModels = append(questionModels, mapToModel(question, userID, note.ID))
	}

	if err := m.repository.CreateMany(ctx, questionModels); err != nil {
		return nil, err
	}

	return &dto.GenerateQuestionsResponse{
		NoteID:    note.ID,
		Questions: questions.Questions,
		Message:   "questions generated successfully",
	}, nil
}

func mapToModel(question llm.GeneratedMCQ, userID string, noteID uint) model.Question {
	options := question.Options

	var correctOption string
	for i, option := range options {
		if strings.TrimSpace(option) != strings.TrimSpace(question.Answer) {
			continue
		}

		switch i {
		case 0:
			correctOption = "A"
		case 1:
			correctOption = "B"
		case 2:
			correctOption = "C"
		case 3:
			correctOption = "D"
		}
	}

	return model.Question{
		UserID:        userID,
		NoteID:        noteID,
		QuestionText:  question.Question,
		OptionA:       options[0],
		OptionB:       options[1],
		OptionC:       options[2],
		OptionD:       options[3],
		CorrectOption: correctOption,
		SourceType:    "llm",
	}
}
