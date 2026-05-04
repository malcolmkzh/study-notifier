package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	questionrepository "github.com/malcolmkzh/study-notifier/internal/modules/questions/repository"
	reminderdto "github.com/malcolmkzh/study-notifier/internal/modules/reminder/dto"
	remindermodel "github.com/malcolmkzh/study-notifier/internal/modules/reminder/model"
	reminderrepository "github.com/malcolmkzh/study-notifier/internal/modules/reminder/repository"
	userrepository "github.com/malcolmkzh/study-notifier/internal/modules/user/repository"
	"github.com/malcolmkzh/study-notifier/internal/utilities/errorutil"
	"github.com/malcolmkzh/study-notifier/internal/utilities/notification"
	"github.com/malcolmkzh/study-notifier/internal/utilities/scheduler"
)

const TaskNameSendReminder scheduler.TaskName = "send_reminder"
const TaskNamePlanSmartReminders scheduler.TaskName = "plan_smart_reminders"

const (
	defaultTimezone          = "Asia/Singapore"
	dailySmartReminderTarget = 8
	coldStartReminderTarget  = 24
	minAttemptsForSmartHours = 5
	timingExplorationPercent = 20
)

type ReminderJobMetadata struct {
	ReminderID int64 `json:"reminder_id"`
}

type Implementation struct {
	reminderRepo        reminderrepository.Utility
	jobRepo             scheduler.JobRepository
	questionRepo        questionrepository.Utility
	userRepo            userrepository.Utility
	notificationUtility notification.Utility
	schedulerUtility    scheduler.Utility
}

func NewService(
	reminderRepo reminderrepository.Utility,
	jobRepo scheduler.JobRepository,
	questionRepo questionrepository.Utility,
	userRepo userrepository.Utility,
	notificationUtility notification.Utility,
	schedulerUtility scheduler.Utility,
) (*Implementation, error) {
	if reminderRepo == nil {
		return nil, errors.New("reminder repository is required")
	}
	if jobRepo == nil {
		return nil, errors.New("job repository is required")
	}
	if questionRepo == nil {
		return nil, errors.New("question repository is required")
	}
	if userRepo == nil {
		return nil, errors.New("user repository is required")
	}
	if notificationUtility == nil {
		return nil, errors.New("notification utility is required")
	}
	if schedulerUtility == nil {
		return nil, errors.New("scheduler utility is required")
	}

	return &Implementation{
		reminderRepo:        reminderRepo,
		jobRepo:             jobRepo,
		questionRepo:        questionRepo,
		userRepo:            userRepo,
		notificationUtility: notificationUtility,
		schedulerUtility:    schedulerUtility,
	}, nil
}

func (s *Implementation) CreateReminder(ctx context.Context, req CreateReminderRequest) error {
	if req.UserID == "" {
		return errorutil.NewWithMessage(errorutil.CodeValidation, "user id is required")
	}
	if req.ScheduledAt.IsZero() {
		return errorutil.NewWithMessage(errorutil.CodeValidation, "scheduled at is required")
	}

	account, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return err
	}
	if account == nil {
		return errorutil.NewWithMessage(errorutil.CodeNotFound, "user not found")
	}
	if account.TelegramChatID == nil || strings.TrimSpace(*account.TelegramChatID) == "" {
		return errorutil.New(errorutil.CodeTelegramNotLinked)
	}

	reminder := remindermodel.Reminder{
		UserID:         req.UserID,
		TelegramChatID: strings.TrimSpace(*account.TelegramChatID),
		ScheduledAt:    req.ScheduledAt,
		Status:         remindermodel.ReminderStatusPending,
	}

	if err := s.reminderRepo.Create(ctx, &reminder); err != nil {
		return err
	}

	metadata, err := json.Marshal(ReminderJobMetadata{
		ReminderID: reminder.ID,
	})
	if err != nil {
		return err
	}

	job := scheduler.Job{
		Name:        fmt.Sprintf("send-reminder-%d", reminder.ID),
		TaskName:    TaskNameSendReminder,
		Status:      scheduler.JobStatusActive,
		Metadata:    metadata,
		ScheduledAt: reminder.ScheduledAt,
	}

	if err := s.jobRepo.CreateJob(ctx, &job); err != nil {
		return err
	}

	return s.schedulerUtility.SyncLocalJobWithDBJob(ctx, job)
}

// Reminder Settings Related Service Methods
func (s *Implementation) GetReminderSetting(ctx context.Context, userID string) (*reminderdto.ReminderSettingResponse, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errorutil.NewWithMessage(errorutil.CodeValidation, "user id is required")
	}

	setting, err := s.reminderRepo.GetSettingByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return &reminderdto.ReminderSettingResponse{
			UserID:   userID,
			Enabled:  false,
			Timezone: defaultTimezone,
		}, nil
	}

	return mapSettingResponse(*setting), nil
}

func (s *Implementation) UpdateReminderSetting(ctx context.Context, req UpdateReminderSettingRequest) (*reminderdto.ReminderSettingResponse, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errorutil.NewWithMessage(errorutil.CodeValidation, "user id is required")
	}

	setting := remindermodel.ReminderSetting{
		UserID:   req.UserID,
		Enabled:  req.Enabled,
		Timezone: defaultTimezone,
	}
	if err := s.reminderRepo.UpsertSetting(ctx, &setting); err != nil {
		return nil, err
	}

	if !req.Enabled {
		if err := s.cancelFutureRemindersForUser(ctx, req.UserID); err != nil {
			return nil, err
		}
	}

	return s.GetReminderSetting(ctx, req.UserID)
}

// Scheduler Job Handler for Daily Smart Reminder Planning
func (s *Implementation) HandlePlannerJob(ctx context.Context, metadata json.RawMessage) error {
	_ = metadata

	if err := s.PlanSmartReminders(ctx, nil); err != nil {
		return err
	}

	return s.EnsureNextPlannerJob(ctx)
}

func (s *Implementation) PlanSmartReminders(ctx context.Context, userID *string) error {
	if userID != nil {
		setting, err := s.reminderRepo.GetSettingByUserID(ctx, *userID)
		if err != nil {
			return err
		}
		if setting == nil || !setting.Enabled {
			return errorutil.NewWithMessage(errorutil.CodeValidation, "smart reminders are not enabled")
		}

		return s.planForSetting(ctx, *setting)
	}

	settings, err := s.reminderRepo.ListEnabledSettings(ctx)
	if err != nil {
		return err
	}

	for _, setting := range settings {
		if err := s.planForSetting(ctx, setting); err != nil {
			// Keep planning other users even if one user is not ready.
			continue
		}
	}

	return nil
}

// Create next daily scheduled job to plan smart reminders
func (s *Implementation) EnsureNextPlannerJob(ctx context.Context) error {
	location, err := time.LoadLocation(defaultTimezone)
	if err != nil {
		return err
	}

	now := time.Now().In(location)
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, location)
	jobName := fmt.Sprintf("plan-smart-reminders-%s", nextMidnight.Format("2006-01-02"))

	jobs, err := s.jobRepo.SelectJobs(ctx, scheduler.SelectJobsRequest{
		Statuses: []scheduler.JobStatus{scheduler.JobStatusActive},
	})
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if job.Name == jobName {
			return nil
		}
	}

	job := scheduler.Job{
		Name:        jobName,
		TaskName:    TaskNamePlanSmartReminders,
		Status:      scheduler.JobStatusActive,
		ScheduledAt: nextMidnight.UTC(),
	}

	if err := s.jobRepo.CreateJob(ctx, &job); err != nil {
		return err
	}

	return s.schedulerUtility.SyncLocalJobWithDBJob(ctx, job)
}

// Scheduler Job Handler for Sending Reminders
func (s *Implementation) HandleSendReminderJob(ctx context.Context, metadata json.RawMessage) error {
	var data ReminderJobMetadata
	if err := json.Unmarshal(metadata, &data); err != nil {
		return err
	}

	return s.TriggerReminder(ctx, data.ReminderID)
}

func (s *Implementation) TriggerReminder(ctx context.Context, reminderID int64) error {
	reminder, err := s.reminderRepo.GetByID(ctx, reminderID)
	if err != nil {
		return err
	}
	if reminder == nil {
		return fmt.Errorf("reminder %d not found", reminderID)
	}

	question, err := s.questionRepo.GetRandomQuestionByUserID(ctx, reminder.UserID)
	if err != nil {
		return err
	}
	if question == nil {
		return fmt.Errorf("no question found for user %s", reminder.UserID)
	}

	correctOptionID, err := correctOptionID(question.CorrectOption)
	if err != nil {
		return err
	}

	attempt := remindermodel.QuestionAttempt{
		UserID:          reminder.UserID,
		ReminderID:      reminder.ID,
		QuestionID:      question.ID,
		CorrectOptionID: correctOptionID,
		SentAt:          time.Now().UTC(),
	}
	if err := s.reminderRepo.CreateQuestionAttempt(ctx, &attempt); err != nil {
		return err
	}

	pollID, err := s.notificationUtility.SendTelegramQuizPoll(ctx, notification.SendTelegramQuizPollRequest{
		ChatID:          reminder.TelegramChatID,
		Question:        question.QuestionText,
		Options:         []string{question.OptionA, question.OptionB, question.OptionC, question.OptionD},
		CorrectOptionID: attempt.CorrectOptionID,
	})
	if err != nil {
		return err
	}
	if strings.TrimSpace(pollID) != "" {
		if err := s.reminderRepo.UpdateQuestionAttemptTelegramPollID(ctx, attempt.ID, pollID); err != nil {
			return err
		}
	}

	return s.reminderRepo.MarkSent(ctx, reminderID)
}

func correctOptionID(correctOption string) (int, error) {
	switch strings.ToUpper(strings.TrimSpace(correctOption)) {
	case "A":
		return 0, nil
	case "B":
		return 1, nil
	case "C":
		return 2, nil
	case "D":
		return 3, nil
	default:
		return 0, fmt.Errorf("invalid correct option %q", correctOption)
	}
}

func (s *Implementation) cancelFutureRemindersForUser(ctx context.Context, userID string) error {
	reminders, err := s.reminderRepo.ListPendingByUserIDAfter(ctx, userID, time.Now().UTC())
	if err != nil {
		return err
	}

	activeStatus := scheduler.JobStatusActive
	deletedStatus := scheduler.JobStatusDeleted
	jobs, err := s.jobRepo.SelectJobs(ctx, scheduler.SelectJobsRequest{
		Statuses: []scheduler.JobStatus{activeStatus},
	})
	if err != nil {
		return err
	}

	for _, reminder := range reminders {
		if err := s.reminderRepo.MarkCancelled(ctx, reminder.ID); err != nil {
			return err
		}

		jobName := fmt.Sprintf("send-reminder-%d", reminder.ID)
		for _, job := range jobs {
			if job.Name != jobName {
				continue
			}

			if err := s.jobRepo.UpdateJob(ctx, scheduler.UpdateJobRequest{
				ID:     job.ID,
				Status: &deletedStatus,
			}); err != nil {
				return err
			}

			dbJob, err := s.jobRepo.GetJobByID(ctx, job.ID)
			if err != nil {
				return err
			}
			if dbJob != nil {
				if err := s.schedulerUtility.SyncLocalJobWithDBJob(ctx, *dbJob); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func mapSettingResponse(setting remindermodel.ReminderSetting) *reminderdto.ReminderSettingResponse {
	return &reminderdto.ReminderSettingResponse{
		UserID:   setting.UserID,
		Enabled:  setting.Enabled,
		Timezone: settingTimezone(setting),
	}
}

func settingTimezone(setting remindermodel.ReminderSetting) string {
	if strings.TrimSpace(setting.Timezone) == "" {
		return defaultTimezone
	}

	return strings.TrimSpace(setting.Timezone)
}
