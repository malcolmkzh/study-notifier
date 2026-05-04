package service

import (
	"context"
	"crypto/rand"
	"math/big"
	"sort"
	"strings"
	"time"

	remindermodel "github.com/malcolmkzh/study-notifier/internal/modules/reminder/model"
)

func (s *Implementation) planForSetting(ctx context.Context, setting remindermodel.ReminderSetting) error {
	account, err := s.userRepo.GetByID(ctx, setting.UserID)
	if err != nil {
		return err
	}
	if account == nil || account.TelegramChatID == nil || strings.TrimSpace(*account.TelegramChatID) == "" {
		return nil
	}

	question, err := s.questionRepo.GetRandomQuestionByUserID(ctx, setting.UserID)
	if err != nil {
		return err
	}
	if question == nil {
		return nil
	}

	location, err := time.LoadLocation(settingTimezone(setting))
	if err != nil {
		location, _ = time.LoadLocation(defaultTimezone)
	}

	now := time.Now().In(location)
	planningStart := now.UTC()
	planningEnd := now.Add(24 * time.Hour).UTC()

	attempts, err := s.reminderRepo.ListAnsweredAttemptsByUserID(ctx, setting.UserID)
	if err != nil {
		return err
	}

	target := dailySmartReminderTarget
	if isColdStart(attempts) {
		target = coldStartReminderTarget
	}

	existing, err := s.reminderRepo.ListPendingByUserIDBetween(ctx, setting.UserID, planningStart, planningEnd)
	if err != nil {
		return err
	}

	remaining := target - len(existing)
	if remaining <= 0 {
		return nil
	}

	slots := pickReminderSlots(now, location, attempts, remaining)

	for _, scheduledAt := range slots {
		if err := s.CreateReminder(ctx, CreateReminderRequest{
			UserID:      setting.UserID,
			ScheduledAt: scheduledAt,
		}); err != nil {
			return err
		}
	}

	return nil
}

func pickReminderSlots(now time.Time, location *time.Location, attempts []remindermodel.QuestionAttempt, count int) []time.Time {
	if isColdStart(attempts) {
		return hourlySlots(now, location, count)
	}

	hours := bestAnsweredHours(attempts, location)
	if len(hours) == 0 {
		return hourlySlots(now, location, count)
	}

	var slots []time.Time
	for i := 0; i < count; i++ {
		hour := hours[i%len(hours)]
		if shouldExploreTiming() {
			hour = exploratoryHour(hours)
		}
		slots = append(slots, nextSlotForHour(now, location, hour))
	}

	return slots
}

func isColdStart(attempts []remindermodel.QuestionAttempt) bool {
	return len(attempts) < minAttemptsForSmartHours
}

func hourlySlots(now time.Time, location *time.Location, count int) []time.Time {
	nextHour := now.In(location).Truncate(time.Hour).Add(time.Hour)

	var slots []time.Time
	for i := 0; i < count; i++ {
		slots = append(slots, nextHour.Add(time.Duration(i)*time.Hour).UTC())
	}

	return slots
}

func nextSlotForHour(now time.Time, location *time.Location, hour int) time.Time {
	localNow := now.In(location)
	slot := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), hour, 0, 0, 0, location)
	if !slot.After(localNow) {
		slot = slot.AddDate(0, 0, 1)
	}

	return slot.UTC()
}

func shouldExploreTiming() bool {
	value, err := secureRandomInt(100)
	if err != nil {
		return false
	}

	return value < timingExplorationPercent
}

func exploratoryHour(bestHours []int) int {
	excluded := make(map[int]bool, len(bestHours))
	for _, hour := range bestHours {
		excluded[hour] = true
	}

	var candidates []int
	for hour := 0; hour < 24; hour++ {
		if !excluded[hour] {
			candidates = append(candidates, hour)
		}
	}

	if len(candidates) == 0 {
		return bestHours[0]
	}

	index, err := secureRandomInt(len(candidates))
	if err != nil {
		return candidates[0]
	}

	return candidates[index]
}

func secureRandomInt(max int) (int, error) {
	value, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}

	return int(value.Int64()), nil
}

func bestAnsweredHours(attempts []remindermodel.QuestionAttempt, location *time.Location) []int {
	scoreByHour := make(map[int]int)
	for index, attempt := range attempts {
		if attempt.AnsweredAt == nil {
			continue
		}

		hour := attempt.AnsweredAt.In(location).Hour()
		score := 1
		if attempt.IsCorrect != nil && *attempt.IsCorrect {
			score++
		}
		if index < 20 {
			score++
		}
		scoreByHour[hour] += score
	}

	type hourScore struct {
		hour  int
		score int
	}

	var scores []hourScore
	for hour, score := range scoreByHour {
		scores = append(scores, hourScore{hour: hour, score: score})
	}

	sort.Slice(scores, func(i, j int) bool {
		if scores[i].score == scores[j].score {
			return scores[i].hour < scores[j].hour
		}
		return scores[i].score > scores[j].score
	})

	var hours []int
	for _, score := range scores {
		hours = append(hours, score.hour)
		if len(hours) == dailySmartReminderTarget {
			break
		}
	}

	return hours
}
