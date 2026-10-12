// Package fetcher downloads HTML pages.
package fetcher

import (
	"context"
	"fmt"
	"io"
	http "net/http"
	"strings"
	"time"
)

// 5 MB
const maxBodyBytes = 5 << 20

type LinkFetcher interface {
	Fetch(ctx context.Context, rawURL string) ([]byte, error)
}

type Fetcher struct {
	client    *http.Client
	userAgent string
}

var _ LinkFetcher = (*Fetcher)(nil)

func NewFetcher(timeout time.Duration, userAgent string) *Fetcher {
	return &Fetcher{
		client:    &http.Client{Timeout: timeout},
		userAgent: userAgent,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, linkURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, linkURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request for %s: %w", linkURL, err)
	}

	req.Header.Set("User-Agent", f.userAgent)

	response, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", linkURL, err)
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: unexpected status %d", linkURL, response.StatusCode)
	}

	ct := response.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		return nil, fmt.Errorf("fetch %s: unsupported content type %q", linkURL, ct)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("read body of %s: %w", linkURL, err)
	}

	return body, nil
}
