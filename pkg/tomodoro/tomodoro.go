package tomodoro

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type Client struct {
	httpBaseUrl string
	httpClient  *http.Client
}

type ErrorResponse struct {
	Href  string `json:"href"`
	Error struct {
		Error   int    `json:"error"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewClient() *Client {
	return &Client{
		httpBaseUrl: BaseURLV1,
		httpClient: &http.Client{
			Timeout: HttpClientTimeout,
		},
	}
}

func (c *Client) sendRequest(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// StatusGone must be checked first, because it is a valid response for Timer.StopTimer
	// This Endpoint returns a none 2xx status code, but it is not an error
	if res.StatusCode != http.StatusGone {
		if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
			return NewRequestError(res)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return errors.New("failed to parse response body")
	}
	return nil
}

func (c *Client) createRequestBody(b interface{}, bb *bytes.Buffer) error {

	mBody, err := json.Marshal(b)
	if err != nil {
		return errors.New("failed to parse request body")
	}

	if _, err := bb.Write(mBody); err != nil {
		return errors.New("failed to write request body")
	}
	return nil
}
