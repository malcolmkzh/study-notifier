package service

import (
	"context"
	"strings"
	"time"

	"github.com/malcolmkzh/study-notifier/internal/utilities/errorutil"
)

func (s *Implementation) HandlePollAnswer(ctx context.Context, pollID string, optionIDs []int) error {
	pollID = strings.TrimSpace(pollID)
	if pollID == "" {
		return errorutil.NewWithMessage(errorutil.CodeValidation, "poll id is required")
	}
	if len(optionIDs) == 0 {
		return nil
	}

	attempt, err := s.reminderRepo.GetQuestionAttemptByTelegramPollID(ctx, pollID)
	if err != nil {
		return err
	}
	if attempt == nil {
		return errorutil.NewWithMessage(errorutil.CodeNotFound, "question attempt not found")
	}
	if attempt.AnsweredAt != nil {
		return nil
	}

	isCorrect := optionIDs[0] == attempt.CorrectOptionID
	return s.reminderRepo.MarkQuestionAttemptPollAnswered(ctx, attempt.ID, time.Now().UTC(), isCorrect)
}
