package process

import (
	"context"
	"sync"

	"siteprobe/internal/crawler"
)

func runConcurrent(
	ctx context.Context,
	c *crawler.Crawler,
	urls []string,
	workers int,
) <-chan crawler.Result {
	workers = min(workers, len(urls))

	jobs := make(chan string, len(urls))
	for _, u := range urls {
		jobs <- u
	}
	close(jobs)

	processed := make(chan crawler.Result)

	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for u := range jobs {
				select {
				case processed <- c.Process(ctx, u):
				case <-ctx.Done():
					return
				}
			}
		})
	}

	go func() {
		wg.Wait()
		close(processed)
	}()

	return processed
}
