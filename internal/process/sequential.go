package process

import (
	"context"

	"siteprobe/internal/crawler"
)

func runSequential(
	ctx context.Context,
	c *crawler.Crawler,
	urls []string,
	show bool,
) (failed int) {
	for _, u := range urls {
		if report(c.Process(ctx, u), show) {
			failed++
		}
	}
	return failed
}
