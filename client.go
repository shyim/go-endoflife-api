// Package endoflife provides a Go client for the endoflife.date API (v1).
//
// endoflife.date documents EOL dates and support lifecycles for various
// products. See https://endoflife.date/docs/api/v1/ for the API reference.
//
// Basic usage:
//
//	client := endoflife.NewClient()
//	product, err := client.Product(context.Background(), "ubuntu")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(product.Result.Label)
package endoflife

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DefaultBaseURL is the production base URL of the endoflife.date v1 API.
const DefaultBaseURL = "https://endoflife.date/api/v1"

// defaultUserAgent is sent with every request unless overridden.
const defaultUserAgent = "go-endoflife-api/1.0 (+https://github.com/shyim/go-endoflife-api)"

// HTTPDoer is the subset of *http.Client used by Client. It allows callers to
// inject custom HTTP clients (e.g. with instrumentation or middleware).
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a client for the endoflife.date API.
//
// A zero Client is not ready for use; create one with NewClient.
type Client struct {
	baseURL    string
	userAgent  string
	httpClient HTTPDoer
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the API base URL. Useful for testing or proxies.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithHTTPClient sets a custom HTTP client (or any HTTPDoer).
func WithHTTPClient(doer HTTPDoer) Option {
	return func(c *Client) {
		if doer != nil {
			c.httpClient = doer
		}
	}
}

// WithUserAgent overrides the User-Agent header sent with each request.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}

// NewClient creates a new endoflife.date API client.
func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:    DefaultBaseURL,
		userAgent:  defaultUserAgent,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// APIError is returned when the API responds with a non-success status code.
type APIError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int
	// Status is the HTTP status line (may be empty per RFC9110).
	Status string
	// Body is the raw response body (the API returns HTML for 404s).
	Body string
	// RetryAfter is the value of the Retry-After header, when present (mainly 429).
	RetryAfter time.Duration
}

// Error implements the error interface.
func (e *APIError) Error() string {
	switch e.StatusCode {
	case http.StatusNotFound:
		return "endoflife: resource not found (404)"
	case http.StatusTooManyRequests:
		if e.RetryAfter > 0 {
			return fmt.Sprintf("endoflife: too many requests (429), retry after %s", e.RetryAfter)
		}
		return "endoflife: too many requests (429)"
	default:
		return fmt.Sprintf("endoflife: unexpected status %d", e.StatusCode)
	}
}

// IsNotFound reports whether err is an APIError with a 404 status code.
func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

// IsTooManyRequests reports whether err is an APIError with a 429 status code.
func IsTooManyRequests(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusTooManyRequests
}

// get performs a GET request against path (relative to the base URL) and
// decodes the JSON response body into out.
func (c *Client) get(ctx context.Context, path string, out any) error {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return fmt.Errorf("endoflife: building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("endoflife: performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		return &APIError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(body),
			RetryAfter: parseRetryAfter(resp.Header.Get("Retry-After")),
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("endoflife: decoding response: %w", err)
	}
	return nil
}

// parseRetryAfter parses a Retry-After header value, supporting both the
// delta-seconds and HTTP-date forms.
func parseRetryAfter(v string) time.Duration {
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
		if secs < 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}
