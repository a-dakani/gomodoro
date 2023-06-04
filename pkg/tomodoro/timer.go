package tomodoro

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"time"
)

type Timer struct {
	Href  string `json:"href"`
	Timer struct {
		Name     string    `json:"name"`
		Duration int64     `json:"duration"`
		Start    time.Time `json:"start"`
	} `json:"timer"`
}

type StartTimerRequest struct {
	Duration int64  `json:"duration"`
	Name     string `json:"name"`
}

type StopTimerResponse struct {
	Href    string `json:"href"`
	Type    string `json:"type"`
	Message string `json:"Message"`
}

func (c *Client) StartTimer(ctx context.Context, teamSlug string, duration int64, name string) (*Timer, error) {
	u, err := url.JoinPath(c.httpBaseUrl, URLTeamSlug, teamSlug, URLTimerSlug, URLStartTimerSlug)
	if err != nil {
		return nil, err
	}

	body := StartTimerRequest{
		Duration: duration,
		Name:     name,
	}

	var bBody bytes.Buffer

	if err := c.createRequestBody(&body, &bBody); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, u, &bBody)
	if err != nil {
		return nil, err
	}

	res := Timer{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) StopTimer(ctx context.Context, teamSlug string) (*StopTimerResponse, error) {
	u, err := url.JoinPath(c.httpBaseUrl, URLTeamSlug, teamSlug, URLTimerSlug)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	res := StopTimerResponse{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
