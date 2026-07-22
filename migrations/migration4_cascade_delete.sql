-- Adds ON DELETE CASCADE to every child of users(uid)

PRAGMA foreign_keys = OFF;

BEGIN TRANSACTION;

CREATE TABLE sessions_new (
    session_id   INTEGER PRIMARY KEY,
    uid          INTEGER REFERENCES users(uid) ON DELETE CASCADE,
    started_at   TIMESTAMP NOT NULL,
    completed_at TIMESTAMP NOT NULL,
    duration     INTEGER,
    local_date   TEXT NOT NULL
);
INSERT INTO sessions_new (session_id, uid, started_at, completed_at, duration, local_date)
    SELECT session_id, uid, started_at, completed_at, duration, local_date FROM sessions;
DROP TABLE sessions;
ALTER TABLE sessions_new RENAME TO sessions;

CREATE TABLE active_sessions_new (
    active_session_id INTEGER PRIMARY KEY,
    uid               INTEGER UNIQUE NOT NULL REFERENCES users(uid) ON DELETE CASCADE,
    started_at        TIMESTAMP NOT NULL,
    elapsed_seconds   INTEGER NOT NULL DEFAULT 0 CHECK (elapsed_seconds >= 0),
    last_updated      TIMESTAMP NOT NULL,
    local_date        TEXT NOT NULL
);
INSERT INTO active_sessions_new (active_session_id, uid, started_at, elapsed_seconds, last_updated, local_date)
    SELECT active_session_id, uid, started_at, elapsed_seconds, last_updated, local_date FROM active_sessions;
DROP TABLE active_sessions;
ALTER TABLE active_sessions_new RENAME TO active_sessions;

CREATE TABLE user_settings_new (
    uid                  INTEGER PRIMARY KEY REFERENCES users(uid) ON DELETE CASCADE,
    theme                TEXT,
    session_length       INTEGER NOT NULL,
    target_session_count INTEGER NOT NULL,
    timezone             TEXT,
    CONSTRAINT max_sessions_time CHECK (session_length * target_session_count < 43200)
);
INSERT INTO user_settings_new (uid, theme, session_length, target_session_count, timezone)
    SELECT uid, theme, session_length, target_session_count, timezone FROM user_settings;
DROP TABLE user_settings;
ALTER TABLE user_settings_new RENAME TO user_settings;

CREATE TABLE tokens_new (
    hash   TEXT UNIQUE NOT NULL,
    uid    INTEGER NOT NULL REFERENCES users(uid) ON DELETE CASCADE,
    expiry TIMESTAMP NOT NULL
);
INSERT INTO tokens_new (hash, uid, expiry)
    SELECT hash, uid, expiry FROM tokens;
DROP TABLE tokens;
ALTER TABLE tokens_new RENAME TO tokens;

PRAGMA foreign_key_check;

COMMIT;

PRAGMA foreign_keys = ON;
