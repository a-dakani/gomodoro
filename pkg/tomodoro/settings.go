package tomodoro

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
)

type Settings struct {
	Focus int64 `json:"focus"`
	Pause int64 `json:"pause"`
}

type UpdateSettingsResponse struct {
	Href     string   `json:"href"`
	Settings Settings `json:"Settings"`
}

func (c *Client) UpdateSettings(ctx context.Context, team string, focus int64, pause int64) (*UpdateSettingsResponse, error) {
	u, err := url.JoinPath(c.BaseUrl, URLTeamSlug, team, URLSettingsSlug)
	if err != nil {
		return nil, err
	}

	body := Settings{
		Focus: focus,
		Pause: pause,
	}

	var bBody bytes.Buffer

	if err := c.createRequestBody(&body, &bBody); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, &bBody)
	if err != nil {
		return nil, err
	}

	res := UpdateSettingsResponse{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
