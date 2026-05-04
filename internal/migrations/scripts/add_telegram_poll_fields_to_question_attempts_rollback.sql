ALTER TABLE question_attempts
    DROP INDEX idx_question_attempts_telegram_poll_id,
    DROP COLUMN correct_option_id,
    DROP COLUMN telegram_poll_id;
