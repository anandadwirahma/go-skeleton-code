package httpgateway

import "time"

// config holds the internal configuration resolved from Option values.
type config struct {
	baseURL       string
	timeout       time.Duration
	retryCount    int
	retryWaitTime time.Duration
	headers       map[string]string
}

// Option mutates the HTTPGateway configuration using the functional options pattern.
type Option func(*config)

// WithBaseURL sets the base URL prepended to every request path.
func WithBaseURL(url string) Option {
	return func(c *config) {
		c.baseURL = url
	}
}

// WithTimeout overrides the default request timeout (30s).
func WithTimeout(d time.Duration) Option {
	return func(c *config) {
		c.timeout = d
	}
}

// WithRetryCount overrides the default retry count (3).
func WithRetryCount(n int) Option {
	return func(c *config) {
		c.retryCount = n
	}
}

// WithRetryWaitTime overrides the wait duration between retries (2s).
func WithRetryWaitTime(d time.Duration) Option {
	return func(c *config) {
		c.retryWaitTime = d
	}
}

// WithHeaders merges the given headers on top of the default Content-Type/Accept headers.
func WithHeaders(headers map[string]string) Option {
	return func(c *config) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}
