package tomodoro

import "time"

const (
	BaseURLV1         = "https://api.tomodoro.de/api/v1/"
	TeamSlug          = "team"
	TimerSlug         = "timer"
	StartTimerSlug    = "start"
	SettingsSlug      = "settings"
	HttpClientTimeout = time.Minute
)
