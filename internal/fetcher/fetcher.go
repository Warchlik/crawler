// Package fetcher
package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	config "crawler/internal/config"
)

// 5 MB
const maxBodyBytes = 5 << 20

type Fetcher struct {
	client *http.Client
}

func NewFetcher(timeout time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{Timeout: timeout},
	}
}

func (f *Fetcher) Fetch(ctx context.Context, cfg *config.Config) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", cfg.Url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request for %s: %w", cfg.Url, err)
	}

	req.Header.Set("User-Agent", "")

	response, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", cfg.Url, err)
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: unexpected status %d", cfg.Url, response.StatusCode)
	}

	ct := response.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		return nil, fmt.Errorf("fetch %s: unsupported content type %q", cfg.Url, ct)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("read body of %s: %w", cfg.Url, err)
	}

	return body, nil
}
