package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	questionrepository "github.com/malcolmkzh/study-notifier/internal/modules/questions/repository"
	remindermodel "github.com/malcolmkzh/study-notifier/internal/modules/reminder/model"
	reminderrepository "github.com/malcolmkzh/study-notifier/internal/modules/reminder/repository"
	"github.com/malcolmkzh/study-notifier/internal/utilities/notification"
	"github.com/malcolmkzh/study-notifier/internal/utilities/scheduler"
)

const TaskNameSendReminder scheduler.TaskName = "send_reminder"

type ReminderJobMetadata struct {
	ReminderID int64 `json:"reminder_id"`
}

type Implementation struct {
	reminderRepo        reminderrepository.Utility
	jobRepo             scheduler.JobRepository
	questionRepo        questionrepository.Utility
	notificationUtility notification.Utility
	schedulerUtility    scheduler.Utility
}

func NewService(
	reminderRepo reminderrepository.Utility,
	jobRepo scheduler.JobRepository,
	questionRepo questionrepository.Utility,
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
		notificationUtility: notificationUtility,
		schedulerUtility:    schedulerUtility,
	}, nil
}

func (s *Implementation) CreateReminder(ctx context.Context, req CreateReminderRequest) error {
	if req.UserID == "" {
		return errors.New("user id is required")
	}
	if req.TelegramChatID == "" {
		return errors.New("telegram chat id is required")
	}
	if req.ScheduledAt.IsZero() {
		return errors.New("scheduled at is required")
	}

	reminder := remindermodel.Reminder{
		UserID:         req.UserID,
		TelegramChatID: req.TelegramChatID,
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

func (s *Implementation) CancelReminder(ctx context.Context, reminderID int64) error {
	reminder, err := s.reminderRepo.GetByID(ctx, reminderID)
	if err != nil {
		return err
	}
	if reminder == nil {
		return fmt.Errorf("reminder %d not found", reminderID)
	}

	if err := s.reminderRepo.MarkCancelled(ctx, reminderID); err != nil {
		return err
	}

	jobs, err := s.jobRepo.SelectJobs(ctx, scheduler.SelectJobsRequest{
		Statuses: []scheduler.JobStatus{
			scheduler.JobStatusPending,
			scheduler.JobStatusActive,
		},
	})
	if err != nil {
		return err
	}

	jobName := fmt.Sprintf("send-reminder-%d", reminderID)
	for _, job := range jobs {
		if job.Name != jobName {
			continue
		}

		deletedStatus := scheduler.JobStatusDeleted
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
			return s.schedulerUtility.SyncLocalJobWithDBJob(ctx, *dbJob)
		}

		break
	}

	return nil
}

func (s *Implementation) HandleScheduledJob(ctx context.Context, metadata json.RawMessage) error {
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

	if err := s.notificationUtility.SendQuestion(ctx, notification.SendQuestionRequest{
		TelegramChatID: reminder.TelegramChatID,
		QuestionText:   question.QuestionText,
		Options:        []string{question.OptionA, question.OptionB, question.OptionC, question.OptionD},
		CorrectOption:  question.CorrectOption,
	}); err != nil {
		return err
	}

	return s.reminderRepo.MarkSent(ctx, reminderID)
}
