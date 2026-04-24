package octopus

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrAuthRequired is returned when an authenticated endpoint is called without an API key.
var ErrAuthRequired = errors.New("octopus: API key required for this endpoint")

// ErrPageLimitExceeded is returned when pagination reaches the maximum page limit,
// indicating that the returned results are incomplete.
var ErrPageLimitExceeded = errors.New("octopus: page limit exceeded, results may be incomplete")

// APIError represents a non-2xx response from the Octopus Energy API.
type APIError struct {
	StatusCode int
	Body       string
	Message    string // parsed from the API's "detail" JSON field, if present
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("octopus: API %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("octopus: API returned status %d: %s", e.StatusCode, e.Body)
}

// newAPIError creates an APIError, attempting to parse the "detail" field from the response body.
func newAPIError(statusCode int, body string) *APIError {
	e := &APIError{StatusCode: statusCode, Body: body}
	var parsed struct {
		Detail string `json:"detail"`
	}
	if json.Unmarshal([]byte(body), &parsed) == nil && parsed.Detail != "" {
		e.Message = parsed.Detail
	}
	return e
}

// IsNotFound reports whether err is an API error with status 404.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}
