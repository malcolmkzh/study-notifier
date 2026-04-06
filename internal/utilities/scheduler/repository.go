package scheduler

import (
	"context"
	"errors"
	"time"

	"github.com/malcolmkzh/study-notifier/internal/utilities/db"
	"gorm.io/gorm"
)

type JobRepositoryImplementation struct {
	db db.Utility
}

func NewJobRepository(dbUtility db.Utility) (*JobRepositoryImplementation, error) {
	if dbUtility == nil {
		return nil, errors.New("db utility is required")
	}

	return &JobRepositoryImplementation{
		db: dbUtility,
	}, nil
}

func (r *JobRepositoryImplementation) SelectJobs(ctx context.Context, request SelectJobsRequest) ([]Job, error) {
	var jobs []Job

	query := r.db.DB().WithContext(ctx).Model(&Job{})
	if len(request.Statuses) > 0 {
		query = query.Where("status IN ?", request.Statuses)
	}

	if err := query.Order("scheduled_at ASC").Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *JobRepositoryImplementation) CreateJob(ctx context.Context, job *Job) error {
	if job == nil {
		return errors.New("job is required")
	}

	return r.db.DB().WithContext(ctx).Create(job).Error
}

func (r *JobRepositoryImplementation) UpdateJob(ctx context.Context, request UpdateJobRequest) error {
	if request.ID == 0 {
		return errors.New("job id is required")
	}

	updates := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}

	if request.Status != nil {
		updates["status"] = *request.Status
	}
	if request.ScheduledAt != nil {
		updates["scheduled_at"] = *request.ScheduledAt
	}
	if request.ExecutedAt != nil {
		updates["executed_at"] = *request.ExecutedAt
	}

	return r.db.DB().WithContext(ctx).
		Model(&Job{}).
		Where("id = ?", request.ID).
		Updates(updates).Error
}

func (r *JobRepositoryImplementation) GetJobByID(ctx context.Context, id int64) (*Job, error) {
	var job Job

	if err := r.db.DB().WithContext(ctx).First(&job, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &job, nil
}
