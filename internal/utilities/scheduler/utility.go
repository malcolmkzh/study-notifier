package scheduler

import (
	"context"
	"encoding/json"
)

type TaskFunc func(ctx context.Context, metadata json.RawMessage) error

type Repository interface {
	SelectJobs(ctx context.Context, request SelectJobsRequest) ([]Job, error)
	UpdateJob(ctx context.Context, request UpdateJobRequest) error
}

type JobRepository interface {
	Repository
	CreateJob(ctx context.Context, job *Job) error
	GetJobByID(ctx context.Context, id int64) (*Job, error)
}

type Utility interface {
	RegisterTask(taskName TaskName, taskFn TaskFunc)
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	SyncLocalJobWithDBJob(ctx context.Context, job Job) error
	SyncLocalJobsWithDBJobs(ctx context.Context, jobs []Job) error
}
