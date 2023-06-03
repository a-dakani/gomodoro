package tomodoro

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ErrType string

type RequestError struct {
	StatusCode int
	Href       string
	Err        error
}

func NewRequestError(res *http.Response) error {
	resBody := ErrorResponse{}

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return fmt.Errorf("while calling %s got status: %d but failed to parse error response body", res.Request.URL, res.StatusCode)
	}

	return &RequestError{
		StatusCode: res.StatusCode,
		Href:       res.Request.URL.String(),
		Err:        errors.New(resBody.Error.Message),
	}
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("While calling %s got status: %d and error: %v", r.Href, r.StatusCode, r.Err)
}
func (r *RequestError) NotFound() bool {
	return r.StatusCode == http.StatusNotFound
}
func (r *RequestError) BadRequest() bool {
	return r.StatusCode == http.StatusBadRequest
}
func (r *RequestError) Gone() bool {
	return r.StatusCode == http.StatusGone
}
func (r *RequestError) InternalServerError() bool {
	return r.StatusCode == http.StatusInternalServerError
}
