package tomodoro

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
)

// Team is a team received from the API
type Team struct {
	Name     string   `json:"name"`
	Slug     string   `json:"slug"`
	Settings Settings `json:"settings"`
	Href     string   `json:"href"`
	Links    []Link   `json:"links"`
}

type createTeamRequest struct {
	Team string `json:"team"`
}

// CreateTeamResponse is the response received when creating a team
type CreateTeamResponse struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Link is a link in the response
type Link struct {
	Link string `json:"Link"`
	Rel  string `json:"rel"`
	Type string `json:"type"`
}

// GetTeam gets a team
func (c *Client) GetTeam(ctx context.Context, teamSlug string) (*Team, error) {
	u, err := url.JoinPath(c.httpBaseURL, urlTeamSlug, teamSlug)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	res := Team{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CreateTeam creates a team
func (c *Client) CreateTeam(ctx context.Context, teamName string) (*CreateTeamResponse, error) {
	u, err := url.JoinPath(c.httpBaseURL, urlTeamSlug)
	if err != nil {
		return nil, err
	}

	body := createTeamRequest{
		Team: teamName,
	}

	var bBody bytes.Buffer

	if err := c.createRequestBody(&body, &bBody); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, &bBody)
	if err != nil {
		return nil, err
	}

	res := CreateTeamResponse{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
