package process

import (
	"context"
	"sync"

	"siteprobe/internal/crawler"
)

func runSequential(
	ctx context.Context,
	c *crawler.Crawler,
	urls []string,
) <-chan crawler.Result {
	processed := make(chan crawler.Result)

	var wg sync.WaitGroup

	wg.Go(func() {
		defer close(processed)

		for _, u := range urls {
			select {
			case processed <- c.Process(ctx, u):
			case <-ctx.Done():
				return
			}
		}
	})

	return processed
}
