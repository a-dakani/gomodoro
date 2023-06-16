package tomodoro

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type errorResponse struct {
	Href  string `json:"href"`
	Error struct {
		Error   int    `json:"error"`
		Message string `json:"message"`
	} `json:"error"`
}

type requestError struct {
	StatusCode int
	Href       string
	Err        error
}

func newRequestError(res *http.Response) error {
	resBody := errorResponse{}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return fmt.Errorf(
			"while calling %s got status: %d but failed to parse error response body",
			res.Request.URL, res.StatusCode)
	}

	return &requestError{
		StatusCode: res.StatusCode,
		Href:       res.Request.URL.String(),
		Err:        errors.New(resBody.Error.Message),
	}
}

func (r *requestError) Error() string {
	return fmt.Sprintf("While calling %s got status: %d and error: %v", r.Href, r.StatusCode, r.Err)
}
func (r *requestError) NotFound() bool {
	return r.StatusCode == http.StatusNotFound
}
func (r *requestError) BadRequest() bool {
	return r.StatusCode == http.StatusBadRequest
}
func (r *requestError) Gone() bool {
	return r.StatusCode == http.StatusGone
}
func (r *requestError) InternalServerError() bool {
	return r.StatusCode == http.StatusInternalServerError
}
