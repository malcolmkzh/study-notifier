CREATE TABLE IF NOT EXISTS question_attempts (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    reminder_id BIGINT NOT NULL,
    question_id BIGINT UNSIGNED NOT NULL,
    sent_at TIMESTAMP NOT NULL,
    selected_option CHAR(1) NULL,
    answered_at TIMESTAMP NULL,
    is_correct BOOLEAN NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_question_attempts_user_id (user_id),
    KEY idx_question_attempts_reminder_id (reminder_id),
    KEY idx_question_attempts_question_id (question_id),
    KEY idx_question_attempts_answered_at (answered_at)
);
