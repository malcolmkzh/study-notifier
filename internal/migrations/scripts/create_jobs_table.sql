CREATE TABLE IF NOT EXISTS jobs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    task_name VARCHAR(255) NOT NULL,
    status SMALLINT NOT NULL,
    metadata JSON NULL,
    scheduled_at TIMESTAMP NOT NULL,
    executed_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_jobs_name (name),
    KEY idx_jobs_task_name (task_name),
    KEY idx_jobs_status (status),
    KEY idx_jobs_scheduled_at (scheduled_at)
);
