package tomodoro

import "time"

const (
	baseURLV1         = "https://api.tomodoro.de/api/v1/"
	baseWSURLV1       = "wss://api.tomodoro.de/api/v1/"
	urlTeamSlug       = "team"
	urlTimerSlug      = "timer"
	urlStartTimerSlug = "start"
	urlSettingsSlug   = "settings"
	httpClientTimeout = time.Minute
)
