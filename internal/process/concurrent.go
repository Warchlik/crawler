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
) []crawler.Result {
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
				processed <- c.Process(ctx, u)
			}
		})
	}

	go func() {
		wg.Wait()
		close(processed)
	}()

	var results []crawler.Result

	for p := range processed {
		results = append(results, p)
	}

	return results
}
