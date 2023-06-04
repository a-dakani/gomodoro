package tomodoro

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
)

type Team struct {
	Name     string   `json:"name"`
	Slug     string   `json:"slug"`
	Settings Settings `json:"settings"`
	Href     string   `json:"href"`
	Links    []Link   `json:"links"`
}

type CreateTeamRequest struct {
	Team string `json:"team"`
}

type CreateTeamResponse struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Link struct {
	Link string `json:"Link"`
	Rel  string `json:"rel"`
	Type string `json:"type"`
}

func (c *Client) GetTeam(ctx context.Context, teamSlug string) (*Team, error) {
	u, err := url.JoinPath(c.BaseUrl, URLTeamSlug, teamSlug)
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

func (c *Client) CreateTeam(ctx context.Context, teamName string) (*CreateTeamResponse, error) {
	u, err := url.JoinPath(c.BaseUrl, URLTeamSlug)
	if err != nil {
		return nil, err
	}

	body := CreateTeamRequest{
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
