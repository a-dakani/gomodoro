package main

import (
	"context"
	"github.com/a-dakani/gomodoro/pkg/tomodoro"
	tea "github.com/charmbracelet/bubbletea"
)

var ctx = context.Background()
var tc = tomodoro.NewClient()

func getTeam(teamName string) tea.Cmd {
	return func() tea.Msg {
		t, err := tc.GetTeam(ctx, teamName)
		if err != nil {
			return ErrorMsg(err)
		}

		return Team{
			Name:  t.Name,
			Slug:  t.Slug,
			Focus: t.Settings.Focus,
			Pause: t.Settings.Pause,
		}
	}
}

func createTeam(teamName string) tea.Cmd {
	return func() tea.Msg {
		nt, err := tc.CreateTeam(ctx, teamName)
		if err != nil {
			return ErrorMsg(err)
		}
		t, err := tc.GetTeam(ctx, nt.Slug)
		if err != nil {
			return ErrorMsg(err)
		}
		return Team{
			Name:  t.Name,
			Slug:  t.Slug,
			Focus: t.Settings.Focus,
			Pause: t.Settings.Pause,
		}
	}
}
