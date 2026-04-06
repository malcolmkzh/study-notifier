CREATE TABLE IF NOT EXISTS reminders (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    telegram_chat_id VARCHAR(64) NOT NULL,
    scheduled_at TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_reminders_user_id (user_id),
    KEY idx_reminders_status (status),
    KEY idx_reminders_scheduled_at (scheduled_at)
);
