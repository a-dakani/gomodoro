package main

import (
	"encoding/json"
	"errors"
	"os"
)

func readTeamsFile() ([]Team, error) {
	var teams []Team
	if _, err := os.Stat(teamsFile); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(teamsFile)
		if err != nil {
			return teams, err
		}
		team := addTestTeam()
		teams = append(teams, team)
		bytes, err := json.MarshalIndent(teams, "", "  ")
		if err != nil {
			return teams, err
		}
		err = os.WriteFile(teamsFile, bytes, 0644)
		return teams, err
	}
	bytes, err := os.ReadFile(teamsFile)
	if err != nil {
		return teams, err
	}
	err = json.Unmarshal(bytes, &teams)
	if err != nil {
		return teams, err
	}
	return teams, nil
}

func addTestTeam() Team {
	return Team{
		Name: "gomodoro-test-team",
		Slug: "gomodoro-test-team",
	}
}
