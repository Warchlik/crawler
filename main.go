package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	config "siteprobe/internal/config"
	crawler "siteprobe/internal/crawler"
	fetcher "siteprobe/internal/fetcher"
	parser "siteprobe/internal/parser"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	c := crawler.NewCrawler(
		fetcher.NewFetcher(requestTimeout, cfg.UserAgent),
		parser.NewHTMLParser(),
	)

	start := time.Now()
	failed := 0

	for _, u := range cfg.Url {
		res := c.Process(ctx, u)
		if res.Err != nil {
			failed++
			fmt.Printf("FAIL %s (%s): %v\n", res.URL, res.Took.Round(time.Millisecond), res.Err)
			continue
		}
		fmt.Printf("OK %s (%s): %d links\n", res.URL, res.Took.Round(time.Millisecond), len(res.Links))

		if cfg.Show {
			for _, link := range res.Links {
				// Print all links with tab
				fmt.Println("    ", link)
			}
		}
	}

	fmt.Printf("done: %d urls, %d failed, %s\n", len(cfg.Url), failed, time.Since(start).Round(time.Millisecond))
	return nil
}
