package biznetgio

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int
	Code       int    // envelope `code` when parseable
	Message    string // envelope/detail/body text
	Body       string // raw response body, always retained
}

func (e *APIError) Error() string {
	return fmt.Sprintf("biznetgio api error: status %d, code %d, message %s", e.StatusCode, e.Code, e.Message)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}
// wip 1
