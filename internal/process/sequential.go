package process

import (
	"context"

	"siteprobe/internal/crawler"
)

func runSequential(
	ctx context.Context,
	c *crawler.Crawler,
	urls []string,
) <-chan crawler.Result {
	processed := make(chan crawler.Result)

	go func() {
		defer close(processed)

		for _, u := range urls {
			select {
			case processed <- c.Process(ctx, u):
			case <-ctx.Done():
				return
			}
		}
	}()

	return processed
}
