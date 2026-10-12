package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"crawler/internal/config"
	"crawler/internal/fetcher"
	"crawler/internal/parser"
)

const requestTimeout = 10 * time.Second

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	baseURL, err := url.Parse(cfg.Url)
	if err != nil {
		return err
	}

	if baseURL.Host == "" {
		return fmt.Errorf("invalid url %q: use a full address, e.g. https://example.com", cfg.Url)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var (
		f fetcher.LinkFetcher  = fetcher.NewFetcher(requestTimeout, cfg.UserAgent)
		p parser.LinkExtractor = parser.NewHTMLParser()
	)

	body, err := f.Fetch(ctx, cfg.Url)
	if err != nil {
		return err
	}

	links, err := p.Extract(baseURL, body)
	if err != nil {
		return err
	}

	fmt.Printf("found %d links on %s\n", len(links), cfg.Url)
	for _, l := range links {
		fmt.Println(" ", l)
	}

	return nil
}
