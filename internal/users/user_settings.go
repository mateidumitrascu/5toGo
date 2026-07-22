package users

import "time"

type UserSettings struct {
	UserID             int64
	Theme              string
	SessionLength      int
	TargetSessionCount int
	Timezone           *time.Location
}

func NewUserSettings(uid int64, theme string, sessionLength int, targetCount int, timezone *time.Location) *UserSettings {
	return &UserSettings{uid, theme, sessionLength, targetCount, timezone}
}
