package funcs

import (
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/PuerkitoBio/goquery"
)

type crawlTask struct {
	URL   string
	Depth int
}

func extractLinks(page Page) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page.Body))
	if err != nil {
		return nil
	}

	var links []string
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			links = append(links, href)
		}
	})
	return links
}

func normalize(base, raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	if u.Scheme != "" && u.Host != "" {
		return u.String()
	}

	b, err := url.Parse(base)
	if err != nil {
		return ""
	}

	return b.ResolveReference(u).String()
}

func cleanURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	u.Fragment = ""
	return u.String()
}

func sameDomain(a, b string) bool {
	u1, _ := url.Parse(a)
	u2, _ := url.Parse(b)
	return u1.Host == u2.Host
}

func crawlCollect(opts Options, root string) []Page {
	client := defaultClient()

	var (
		taskQueue = make(chan crawlTask, opts.Threads*4)
		results   = make(chan Page, opts.Threads*4)
		visited   sync.Map
		wg        sync.WaitGroup
		jobCount  int64
	)

	// Increment job count for the root
	atomic.AddInt64(&jobCount, 1)

	// Workers
	for i := 0; i < opts.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range taskQueue {

				// Strict depth check
				if task.Depth > opts.Depth {
					atomic.AddInt64(&jobCount, -1)
					continue
				}

				// Dedup
				if _, seen := visited.LoadOrStore(task.URL, true); seen {
					atomic.AddInt64(&jobCount, -1)
					continue
				}

				page, err := fetchURL(client, "GET", task.URL, "", nil)
				if err == nil {
					results <- page
				}

				// Crawl child links
				if opts.Crawl && err == nil {
					for _, link := range extractLinks(page) {
						abs := normalize(task.URL, link)
						if abs == "" {
							continue
						}
						abs = cleanURL(abs)

						if !sameDomain(task.URL, abs) {
							continue
						}

						if _, seen := visited.Load(abs); seen {
							continue
						}

						// Enqueue new job
						atomic.AddInt64(&jobCount, 1)
						taskQueue <- crawlTask{URL: abs, Depth: task.Depth + 1}
					}
				}

				// Job finished
				if atomic.AddInt64(&jobCount, -1) == 0 {
					// If counter hits zero, close queue
					close(taskQueue)
				}
			}
		}()
	}

	// Seed root task
	taskQueue <- crawlTask{URL: root, Depth: 0}

	// Wait for all workers
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect output
	var pages []Page
	for pg := range results {
		pages = append(pages, pg)
	}

	return pages
}
