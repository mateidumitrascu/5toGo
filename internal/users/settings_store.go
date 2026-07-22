package users

import (
	"database/sql"
	"fmt"
	"time"
)

type UserSettingsSQLiteRepo struct {
	db *sql.DB
}

const userSettingsTable = "user_settings"

func NewUserSettingsSQLiteRepo(db *sql.DB) *UserSettingsSQLiteRepo {
	return &UserSettingsSQLiteRepo{db}
}

func (repo *UserSettingsSQLiteRepo) Create(settings *UserSettings) (*UserSettings, error) {
	_, err := repo.db.Exec("INSERT INTO "+userSettingsTable+" (uid, theme, session_length, target_session_count, timezone) VALUES (?, ?, ?, ?, ?)",
		settings.UserID, settings.Theme, settings.SessionLength, settings.TargetSessionCount, settings.Timezone.String())
	if err != nil {
		return nil, fmt.Errorf("repo error creating user settings: %w", err)
	}
	return settings, nil
}

func (repo *UserSettingsSQLiteRepo) Update(settings *UserSettings) (*UserSettings, error) {
	_, err := repo.db.Exec("UPDATE "+userSettingsTable+" SET theme = ?, session_length = ?, target_session_count = ?, timezone = ? WHERE uid = ?",
		settings.Theme, settings.SessionLength, settings.TargetSessionCount, settings.Timezone.String(), settings.UserID)
	if err != nil {
		return nil, fmt.Errorf("repo error creating user settings: %w", err)
	}
	return settings, nil
}

func (repo *UserSettingsSQLiteRepo) Find(uid int64) (*UserSettings, error) {
	row := repo.db.QueryRow("SELECT theme, session_length, target_session_count, timezone FROM "+userSettingsTable+" WHERE uid = ?", uid)

	us := &UserSettings{}
	var tzString string

	err := row.Scan(&us.Theme, &us.SessionLength, &us.TargetSessionCount, &tzString)
	if err != nil {
		return nil, fmt.Errorf("repo error fetching user settings: %w", err)
	}

	us.UserID = uid

	timezone, err := time.LoadLocation(tzString)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone found in db for uid=%d -> %s: %w", uid, tzString, err)
	}
	us.Timezone = timezone

	return us, nil
}
