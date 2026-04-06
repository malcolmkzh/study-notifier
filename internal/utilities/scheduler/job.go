package scheduler

import (
	"encoding/json"
	"time"
)

type JobStatus int

const (
	JobStatusPending JobStatus = iota
	JobStatusActive
	JobStatusFinished
	JobStatusDeleted
)

type TaskName string

type Job struct {
	ID          int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string          `gorm:"column:name;type:varchar(255);uniqueIndex;not null" json:"name"`
	TaskName    TaskName        `gorm:"column:task_name;type:varchar(255);index;not null" json:"task_name"`
	Status      JobStatus       `gorm:"column:status;type:smallint;index;not null" json:"status"`
	Metadata    json.RawMessage `gorm:"column:metadata;type:json" json:"metadata"`
	ScheduledAt time.Time       `gorm:"column:scheduled_at;type:timestamp;not null" json:"scheduled_at"`
	ExecutedAt  *time.Time      `gorm:"column:executed_at;type:timestamp" json:"executed_at,omitempty"`
	UpdatedAt   *time.Time      `gorm:"column:updated_at" json:"updated_at,omitempty"`
	CreatedAt   time.Time       `gorm:"column:created_at;not null" json:"created_at"`
}

func (Job) TableName() string {
	return "jobs"
}

type SelectJobsRequest struct {
	Statuses []JobStatus
}

type UpdateJobRequest struct {
	ID          int64
	Status      *JobStatus
	ScheduledAt *time.Time
	ExecutedAt  *time.Time
}
