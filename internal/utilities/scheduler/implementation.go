package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Implementation struct {
	logger       *slog.Logger
	jobRepo      Repository
	scheduler    gocron.Scheduler
	taskRegistry map[TaskName]TaskFunc
	jobMapper    jobMapper
}

func NewUtility(jobRepo Repository) (*Implementation, error) {
	if jobRepo == nil {
		return nil, fmt.Errorf("job repository is required")
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("create scheduler: %w", err)
	}

	return &Implementation{
		logger:       slog.Default(),
		jobRepo:      jobRepo,
		scheduler:    scheduler,
		taskRegistry: make(map[TaskName]TaskFunc),
		jobMapper:    newJobMapper(),
	}, nil
}

func (s *Implementation) RegisterTask(taskName TaskName, taskFn TaskFunc) {
	s.taskRegistry[taskName] = taskFn
}

func (s *Implementation) Start(ctx context.Context) error {
	activeJobs, err := s.jobRepo.SelectJobs(ctx, SelectJobsRequest{
		Statuses: []JobStatus{JobStatusActive},
	})
	if err != nil {
		return fmt.Errorf("select active jobs: %w", err)
	}

	for _, job := range activeJobs {
		if err := s.SyncLocalJobWithDBJob(ctx, job); err != nil {
			s.logger.ErrorContext(ctx, "failed to load job into local scheduler", "job_name", job.Name, "job_id", job.ID, "error", err)
		}
	}

	s.scheduler.Start()
	return nil
}

func (s *Implementation) Shutdown(ctx context.Context) error {
	s.logger.InfoContext(ctx, "shutting down scheduler")

	if err := s.scheduler.Shutdown(); err != nil {
		return fmt.Errorf("shutdown scheduler: %w", err)
	}

	return nil
}

func (s *Implementation) SyncLocalJobsWithDBJobs(ctx context.Context, jobs []Job) error {
	for _, job := range jobs {
		if err := s.SyncLocalJobWithDBJob(ctx, job); err != nil {
			return fmt.Errorf("sync job %q: %w", job.Name, err)
		}
	}

	return nil
}

func (s *Implementation) SyncLocalJobWithDBJob(ctx context.Context, job Job) error {
	metadata, err := normalizeMetadata(job.Metadata)
	if err != nil {
		return fmt.Errorf("invalid metadata for job %q: %w", job.Name, err)
	}

	taskFunc, ok := s.taskRegistry[job.TaskName]
	if !ok {
		return fmt.Errorf("task %q is not registered", job.TaskName)
	}

	runAt := job.ScheduledAt.UTC()
	now := time.Now().UTC()
	if runAt.Before(now) {
		runAt = now.Add(1 * time.Second)
	}

	definition := gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(runAt))
	task := gocron.NewTask(taskFunc, metadata)
	options := []gocron.JobOption{
		gocron.WithName(job.Name),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				s.logger.Info("starting scheduled job", "job_name", jobName, "scheduler_job_id", jobID.String())
			}),
			gocron.AfterJobRuns(s.afterJobRuns),
			gocron.AfterJobRunsWithError(s.afterJobRunsWithError),
		),
	}

	existingSchedulerJobID, found := s.jobMapper.GetSchedulerJobID(job.ID)
	if found {
		schedulerJobID, parseErr := uuid.Parse(existingSchedulerJobID)
		if parseErr != nil {
			return fmt.Errorf("parse scheduler job id %q: %w", existingSchedulerJobID, parseErr)
		}

		switch job.Status {
		case JobStatusActive:
			if _, err := s.scheduler.Update(schedulerJobID, definition, task, options...); err != nil {
				return fmt.Errorf("update scheduled job %q: %w", job.Name, err)
			}
		case JobStatusDeleted:
			if err := s.scheduler.RemoveJob(schedulerJobID); err != nil {
				return fmt.Errorf("remove scheduled job %q: %w", job.Name, err)
			}
			s.jobMapper.Remove(job.ID)
		}

		return nil
	}

	if job.Status != JobStatusActive {
		return nil
	}

	scheduledJob, err := s.scheduler.NewJob(definition, task, options...)
	if err != nil {
		return fmt.Errorf("schedule job %q: %w", job.Name, err)
	}

	s.jobMapper.Add(job.ID, scheduledJob.ID().String())
	return nil
}

func (s *Implementation) afterJobRuns(jobID uuid.UUID, jobName string) {
	ctx := context.Background()

	dbJobID, found := s.jobMapper.GetDBJobID(jobID.String())
	if !found {
		s.logger.Warn("scheduled job finished but local mapping was missing", "job_name", jobName, "scheduler_job_id", jobID.String())
		return
	}

	now := time.Now().UTC()
	finishedStatus := JobStatusFinished
	if err := s.jobRepo.UpdateJob(ctx, UpdateJobRequest{
		ID:         dbJobID,
		Status:     &finishedStatus,
		ExecutedAt: &now,
	}); err != nil {
		s.logger.Error("failed to mark scheduled job as finished", "job_name", jobName, "job_id", dbJobID, "error", err)
		return
	}

	s.jobMapper.Remove(dbJobID)
	s.logger.Info("scheduled job finished successfully", "job_name", jobName, "job_id", dbJobID, "scheduler_job_id", jobID.String())
}

func (s *Implementation) afterJobRunsWithError(jobID uuid.UUID, jobName string, err error) {
	s.logger.Error("scheduled job failed", "job_name", jobName, "scheduler_job_id", jobID.String(), "error", err)
}

func normalizeMetadata(metadata json.RawMessage) (json.RawMessage, error) {
	if len(metadata) == 0 {
		return json.RawMessage(`{}`), nil
	}

	if !json.Valid(metadata) {
		return nil, fmt.Errorf("metadata must be valid JSON")
	}

	return metadata, nil
}
