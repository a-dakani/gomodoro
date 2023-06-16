package tomodoro

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Client is the http Client
type Client struct {
	httpBaseURL string
	httpClient  *http.Client
}

// NewClient creates a new Http Client
func NewClient() *Client {
	return &Client{
		httpBaseURL: baseURLV1,
		httpClient: &http.Client{
			Timeout: httpClientTimeout,
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
			return newRequestError(res)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return errors.New("failed to parse response body")
	}

	return nil
}

func (c *Client) createRequestBody(b interface{}, bb io.Writer) error {
	mBody, err := json.Marshal(b)
	if err != nil {
		return errors.New("failed to parse request body")
	}

	if _, err := bb.Write(mBody); err != nil {
		return errors.New("failed to write request body")
	}

	return nil
}
