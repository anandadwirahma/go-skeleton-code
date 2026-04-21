package httpgateway

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// Response carries the raw outcome of an HTTP call. Callers that supply a
// result pointer to the method functions do not need to inspect Body, but it
// is available for debugging and non-JSON payloads.
type Response struct {
	StatusCode int
	Status     string
	Body       []byte
	Headers    http.Header
}

// HTTPError is returned when the transport succeeds but the server responds
// with a non-2xx status code. Callers can type-assert (or errors.As) on this
// type to branch on StatusCode (e.g. http.StatusNotFound, http.StatusInternalServerError).
type HTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
	URL        string
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("http %d %s: %s (url=%s)", e.StatusCode, e.Status, string(e.Body), e.URL)
}

// AsHTTPError returns the underlying *HTTPError when err wraps one.
// It is a thin convenience wrapper around errors.As for call sites that
// want to inspect the status code without importing "errors".
func AsHTTPError(err error) (*HTTPError, bool) {
	var he *HTTPError
	if errors.As(err, &he) {
		return he, true
	}
	return nil, false
}

// handleResponse normalises a resty response into a Response and returns an
// *HTTPError when the server reported a non-2xx status code.
func handleResponse(resp *resty.Response) (*Response, error) {
	r := &Response{
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Body:       resp.Body(),
		Headers:    resp.Header(),
	}

	if resp.IsError() {
		return r, &HTTPError{
			StatusCode: resp.StatusCode(),
			Status:     resp.Status(),
			Body:       resp.Body(),
			URL:        resp.Request.URL,
		}
	}

	return r, nil
}
