ALTER TABLE question_attempts
    ADD COLUMN telegram_poll_id VARCHAR(128) NULL AFTER question_id,
    ADD COLUMN correct_option_id TINYINT NOT NULL DEFAULT 0 AFTER telegram_poll_id,
    ADD UNIQUE KEY idx_question_attempts_telegram_poll_id (telegram_poll_id);
