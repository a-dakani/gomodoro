package main

import (
	"context"
	"github.com/a-dakani/gomodoro/pkg/tomodoro"
)

var ctx = context.Background()
var tc = tomodoro.NewClient()

func getTeam(teamName string) (Team, error) {
	t, err := tc.GetTeam(ctx, teamName)
	if err != nil {
		return Team{}, err
	}

	return Team{
		Name:  t.Name,
		Slug:  t.Slug,
		Focus: t.Settings.Focus,
		Pause: t.Settings.Pause,
	}, nil
}
