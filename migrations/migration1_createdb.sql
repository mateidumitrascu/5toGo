CREATE TABLE users (
    uid INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP
);

CREATE TABLE sessions (
    session_id INTEGER PRIMARY KEY,
    uid INTEGER REFERENCES users(uid),
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    duration INTEGER
);

CREATE TABLE active_sessions (
    active_session_id INTEGER PRIMARY KEY,
    uid INTEGER REFERENCES users(uid) UNIQUE NOT NULL,
    started_at TIMESTAMP NOT NULL,
    elapsed_seconds INTEGER NOT NULL DEFAULT 0 CHECK (elapsed_seconds >=0),
    last_updated TIMESTAMP NOT NULL
);

CREATE TABLE user_settings (
    uid INTEGER PRIMARY KEY REFERENCES users(uid),
    theme TEXT,
    session_length INTEGER NOT NULL,
    target_session_count INTEGER NOT NULL,
    -- a user can only have a target focus time of 12 hours per day, changeable in the future
    CONSTRAINT max_sessions_time CHECK (session_length * target_session_count < 43200)
);
