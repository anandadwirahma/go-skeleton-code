package httpgateway

import (
	"time"

	"github.com/go-resty/resty/v2"
)

// Default configuration applied when no Option overrides it.
const (
	defaultTimeout       = 30 * time.Second
	defaultRetryCount    = 3
	defaultRetryWaitTime = 2 * time.Second
)

// HTTPGateway wraps a resty client and exposes one method per HTTP verb.
// Construct it with New and reuse a single instance across the application.
type HTTPGateway struct {
	client *resty.Client
}

// New builds an HTTPGateway, applying the given options on top of sensible defaults
// (30s timeout, 3 retries with 2s wait, JSON Content-Type and Accept headers).
func New(opts ...Option) *HTTPGateway {
	cfg := &config{
		timeout:       defaultTimeout,
		retryCount:    defaultRetryCount,
		retryWaitTime: defaultRetryWaitTime,
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	client := resty.New().
		SetTimeout(cfg.timeout).
		SetRetryCount(cfg.retryCount).
		SetRetryWaitTime(cfg.retryWaitTime).
		SetHeaders(cfg.headers)

	if cfg.baseURL != "" {
		client.SetBaseURL(cfg.baseURL)
	}

	return &HTTPGateway{client: client}
}

// Client returns the underlying resty client for advanced customisation
// (middleware, transport tuning, etc.) that is not exposed via Options.
func (h *HTTPGateway) Client() *resty.Client {
	return h.client
}
