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

func removeTeamFromFile(team Team) error {
	teams, err := readTeamsFile()
	if err != nil {
		return err
	}
	for i, t := range teams {
		if t.Slug == team.Slug {
			teams = append(teams[:i], teams[i+1:]...)
			break
		}
	}
	bytes, err := json.MarshalIndent(teams, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(teamsFile, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func addTeamToFile(team Team) error {
	teamExists := false
	fileTeams, err := readTeamsFile()
	if err != nil {
		return err
	}

	// check if team already exists if it does replace it with the new one else append it
	for i, t := range fileTeams {
		if t.Slug == team.Slug {
			fileTeams[i] = team
			teamExists = true
			break
		}
	}
	if !teamExists {
		fileTeams = append(fileTeams, team)
	}

	bytes, err := json.MarshalIndent(fileTeams, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(teamsFile, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func addTestTeam() Team {
	return Team{
		Name:  "gomodoro-test-team",
		Slug:  "gomodoro-test-team",
		Focus: 25 * 60 * 1000000000,
		Pause: 5 * 60 * 1000000000,
	}
}
