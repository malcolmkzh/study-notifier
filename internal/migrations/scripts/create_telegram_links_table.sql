CREATE TABLE IF NOT EXISTS telegram_links (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    code VARCHAR(32) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_telegram_links_code (code),
    KEY idx_telegram_links_user_id (user_id),
    KEY idx_telegram_links_expires_at (expires_at)
);
