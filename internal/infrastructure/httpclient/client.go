// Package httpclient provides a reusable HTTP client wrapper for calling external APIs.
// It supports configurable timeouts and logs all outgoing requests via zap.
package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Client wraps net/http.Client with structured logging and convenience methods.
type Client struct {
	baseURL    string
	httpClient *http.Client
	log        *zap.Logger
}

// New returns a new Client with the given base URL and timeout.
func New(baseURL string, timeoutSec int, log *zap.Logger) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
		log: log,
	}
}

// Get sends a GET request to path (relative to baseURL) and decodes the JSON
// response body into dest.
func (c *Client) Get(ctx context.Context, path string, dest interface{}) error {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("httpclient: create request: %w", err)
	}

	return c.do(req, dest)
}

// Post sends a POST request to path with a JSON-encoded body and decodes the
// response body into dest.
func (c *Client) Post(ctx context.Context, path string, body interface{}, dest interface{}) error {
	url := c.baseURL + path

	raw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("httpclient: marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("httpclient: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.do(req, dest)
}

// do executes an HTTP request, logs the outcome, and decodes the response.
func (c *Client) do(req *http.Request, dest interface{}) error {
	start := time.Now()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error("httpclient: request failed",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Error(err),
		)
		return fmt.Errorf("httpclient: request failed: %w", err)
	}
	defer resp.Body.Close()

	latency := time.Since(start)

	c.log.Info("httpclient: request completed",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Int("status", resp.StatusCode),
		zap.Duration("latency", latency),
	)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("httpclient: unexpected status %d for %s %s",
			resp.StatusCode, req.Method, req.URL.String())
	}

	if dest == nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("httpclient: read body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, dest); err != nil {
		return fmt.Errorf("httpclient: unmarshal response: %w", err)
	}

	return nil
}

// ----- Example usage: JSONPlaceholder -----

// JSONPlaceholderPost represents a post from jsonplaceholder.typicode.com.
type JSONPlaceholderPost struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

// FetchExamplePost demonstrates fetching post #1 from the mock API.
// This can be removed or replaced with real external API calls.
func FetchExamplePost(ctx context.Context, client *Client) (*JSONPlaceholderPost, error) {
	var post JSONPlaceholderPost
	if err := client.Get(ctx, "/posts/1", &post); err != nil {
		return nil, err
	}
	return &post, nil
}
