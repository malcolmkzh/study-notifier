ALTER TABLE notes
ADD COLUMN folder_id BIGINT UNSIGNED NULL AFTER user_id,
ADD INDEX idx_notes_folder_id (folder_id);
