ALTER TABLE sessions ADD COLUMN local_date TEXT NOT NULL;
ALTER TABLE active_sessions ADD COLUMN local_date TEXT NOT NULL;

ALTER TABLE sessions DROP COLUMN completed_at;
ALTER TABLE sessions ADD COLUMN completed_at TIMESTAMP NOT NULL;

ALTER TABLE user_settings ADD COLUMN timezone TEXT;
