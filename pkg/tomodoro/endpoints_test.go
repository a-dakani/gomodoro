package tomodoro

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	teamName       = "gomodoro-test-team"
	focus    int64 = 25 * 60 * 1000 * 1000
	pause    int64 = 5 * 60 * 1000 * 1000
)

var teamSlug = ""

func TestCreateTeam(t *testing.T) {
	tc := NewClient()
	ctx := context.Background()

	team, err := tc.CreateTeam(ctx, teamName)
	assert.NotNil(t, team, "team should not be nil")
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, teamName, team.Name, "team name should match")
	// save team slug for later tests
	teamSlug = team.Slug
}

func TestGetTeam(t *testing.T) {
	tc := NewClient()
	ctx := context.Background()

	team, err := tc.GetTeam(ctx, teamSlug)
	assert.NotNil(t, team, "team should not be nil")
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, teamName, team.Name, "team name should match")
	assert.Equal(t, teamSlug, team.Slug, "team slug should match")
}

func TestUpdateSettings(t *testing.T) {
	tc := NewClient()
	ctx := context.Background()

	settings, err := tc.UpdateSettings(ctx, teamName, focus, pause)
	assert.NotNil(t, settings, "team should not be nil")
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, focus, settings.Settings.Focus, "focus should match")
	assert.Equal(t, pause, settings.Settings.Pause, "pause should match")
}

func TestStartTimer(t *testing.T) {
	tc := NewClient()
	ctx := context.Background()

	timer, err := tc.StartTimer(ctx, teamSlug, focus, "Focus")
	assert.NotNil(t, timer, "timer should not be nil")
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, "Fokusphase", timer.Timer.Name, "name should match")
	assert.Equal(t, focus, timer.Timer.Duration, "duration should match")
}

func TestStopTimer(t *testing.T) {
	tc := NewClient()
	ctx := context.Background()

	timer, err := tc.StopTimer(ctx, teamSlug)
	assert.NotNil(t, timer, "timer should not be nil")
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, "timerStopped", timer.Type, "name should match")
	assert.Equal(t, "timer stopped", timer.Message, "duration should match")
}
