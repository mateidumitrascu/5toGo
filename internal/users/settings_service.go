package users

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type SettingsStore interface {
	Create(s *UserSettings) (*UserSettings, error)
	Update(s *UserSettings) (*UserSettings, error)
	Find(uid int64) (*UserSettings, error)
}

type SettingsService struct {
	settingsStore SettingsStore
}

func NewSettingsService(settingsStore SettingsStore) *SettingsService {
	return &SettingsService{settingsStore}
}

func (srv *SettingsService) InitializeSettings(uid int64, timezone string) error {
	theme := os.Getenv("DEFAULT_THEME")
	sessionLength := os.Getenv("DEFAULT_SESSION_LENGTH")
	sessionCount := os.Getenv("DEFAULT_TARGET_SESSION_COUNT")
	tz, err := time.LoadLocation(timezone)
	if err != nil {
		tz = time.UTC
	}

	length, _ := strconv.Atoi(sessionLength)
	count, _ := strconv.Atoi(sessionCount)

	_, err = srv.settingsStore.Create(NewUserSettings(uid, theme, length, count, tz))
	if err != nil {
		return fmt.Errorf("service error initializing user settings: %w", err)
	}
	return nil
}

func (srv *SettingsService) UserTimezone(uid int64) *time.Location {
	// TODO: add optimized timezone fetch function in repo
	s, err := srv.settingsStore.Find(uid)
	if err != nil {
		return nil
	}
	return s.Timezone
}
