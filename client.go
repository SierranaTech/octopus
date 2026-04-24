package octopus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL   = "https://api.octopus.energy/v1"
	userAgent        = "octopus-go/0.1"
	maxResponseBytes = 10 << 20 // 10 MB
)

// Client is an Octopus Energy API client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Octopus Energy API client.
func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) requireAuth() error {
	if c.apiKey == "" {
		return ErrAuthRequired
	}
	return nil
}

func (c *Client) newRequest(ctx context.Context, path string, query url.Values, auth bool) (*http.Request, error) {
	if auth {
		if err := c.requireAuth(); err != nil {
			return nil, err
		}
	}

	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	if auth {
		req.SetBasicAuth(c.apiKey, "")
	}
	return req, nil
}

func (c *Client) do(req *http.Request, v any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newAPIError(resp.StatusCode, string(body))
	}

	return json.Unmarshal(body, v)
}

// get performs a GET request to a path relative to the base URL.
func (c *Client) get(ctx context.Context, path string, query url.Values, auth bool, v any) error {
	req, err := c.newRequest(ctx, path, query, auth)
	if err != nil {
		return err
	}
	return c.do(req, v)
}

// getURL performs a GET request to an absolute URL (used for pagination next links).
func (c *Client) getURL(ctx context.Context, rawURL string, auth bool, v any) error {
	if auth {
		if err := c.requireAuth(); err != nil {
			return err
		}
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("octopus: invalid next URL: %w", err)
	}

	// Validate origin to prevent SSRF via pagination next links.
	baseOrigin, _ := url.Parse(c.baseURL)
	if parsed.Scheme != baseOrigin.Scheme || parsed.Host != baseOrigin.Host {
		return fmt.Errorf("octopus: refusing to follow next URL to different origin: %s", parsed.Host)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	if auth {
		req.SetBasicAuth(c.apiKey, "")
	}
	return c.do(req, v)
}
