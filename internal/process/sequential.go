package process

import (
	"context"

	"siteprobe/internal/crawler"
)

func runSequential(
	ctx context.Context,
	c *crawler.Crawler,
	urls []string,
) []crawler.Result {
	var results []crawler.Result
	for _, u := range urls {
		results = append(results, c.Process(ctx, u))
	}

	return results
}
