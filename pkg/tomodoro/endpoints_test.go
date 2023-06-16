package tomodoro

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	teamName = "gomodoro-test-team"
	focus    = 25 * time.Minute.Nanoseconds()
	pause    = 5 * time.Minute.Nanoseconds()
)

var teamSlug = "gomodoro-test-team"

// Disable this test to avoid creating a new team every time

//func TestCreateTeam(t *testing.T) {
//	tc := NewClient()
//	ctx := context.Background()
//
//	team, err := tc.CreateTeam(ctx, teamName)
//	assert.NotNil(t, team, "team should not be nil")
//	assert.Nil(t, err, "error should be nil")
//	assert.Equal(t, teamName, team.Name, "team name should match")
//	// save team slug for later tests
//	teamSlug = team.Slug
//}

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
	// stop timer if running
	_, _ = tc.StopTimer(ctx, teamSlug)
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
