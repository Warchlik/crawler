// Package config
package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Url       string
	Workers   int
	MaxDepth  int
	Rate      int
	UserAgent string
}

func NewConfig(url string, workers int, maxDepth int, rate int, userAgent string) *Config {
	cfg := Config{
		Url:       url,
		Workers:   workers,
		MaxDepth:  maxDepth,
		Rate:      rate,
		UserAgent: userAgent,
	}

	return &cfg
}

func ParseConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Url, "url", "", "Crawling url")
	flag.IntVar(&cfg.Workers, "workers", 1, "Flag definites how many workers will connect and search in website")
	flag.IntVar(&cfg.MaxDepth, "max-depth", 1, "Flag definites how deep worker can enter int to website definie by links")
	flag.IntVar(&cfg.Rate, "rate", 1, "Flag define rate")
	flag.StringVar(&cfg.UserAgent, "user-agent", "Crawler-Agent/default", "Flag defines custom user agnet")
	flag.Parse()

	if cfg.Url == "" {
		return nil, fmt.Errorf("url is required")
	}

	return &cfg, nil
}
