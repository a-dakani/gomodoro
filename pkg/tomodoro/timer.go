package tomodoro

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"time"
)

// Timer is a timer received from the API
type Timer struct {
	Href  string `json:"href"`
	Timer struct {
		Name     string    `json:"name"`
		Duration int64     `json:"duration"`
		Start    time.Time `json:"start"`
	} `json:"timer"`
}

type startTimerRequest struct {
	Duration int64  `json:"duration"`
	Name     string `json:"name"`
}

// StopTimerResponse is the response received when stopping a timer
type StopTimerResponse struct {
	Href    string `json:"href"`
	Type    string `json:"type"`
	Message string `json:"Message"`
}

// StartTimer starts a timer
func (c *Client) StartTimer(ctx context.Context, teamSlug string, duration int64, name string) (*Timer, error) {
	u, err := url.JoinPath(c.httpBaseURL, urlTeamSlug, teamSlug, urlTimerSlug, urlStartTimerSlug)
	if err != nil {
		return nil, err
	}

	body := startTimerRequest{
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

// StopTimer stops the timer
func (c *Client) StopTimer(ctx context.Context, teamSlug string) (*StopTimerResponse, error) {
	u, err := url.JoinPath(c.httpBaseURL, urlTeamSlug, teamSlug, urlTimerSlug)
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
