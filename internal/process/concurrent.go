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
	show bool,
) (failed int) {
	workers = min(workers, len(urls))

	jobs := make(chan string, len(urls))
	for _, u := range urls {
		jobs <- u
	}
	close(jobs)

	results := make(chan crawler.Result)

	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for u := range jobs {
				results <- c.Process(ctx, u)
			}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if report(res, show) {
			failed++
		}
	}

	failed = 0

	return failed
}
