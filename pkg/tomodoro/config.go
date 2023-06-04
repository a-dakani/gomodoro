package tomodoro

import "time"

const (
	BaseURLV1         = "https://api.tomodoro.de/api/v1/"
	BaseWSURLV1       = "wss://api.tomodoro.de/api/v1/"
	URLTeamSlug       = "team"
	URLTimerSlug      = "timer"
	URLStartTimerSlug = "start"
	URLSettingsSlug   = "settings"
	HttpClientTimeout = time.Minute
)
