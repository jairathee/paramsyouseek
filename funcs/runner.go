package funcs

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

func Run(opts Options) error {
	if !opts.Silent {
		PrintBanner()
	}

	if err := opts.Normalize(); err != nil {
		return err
	}

	client := defaultClient()
	_ = client // currently only used by fetchURL and other helpers

	var pages []Page

	// 1) Raw request file
	if opts.RequestFile != "" {
		method, u, headers, body, err := parseRawRequestFile(opts.RequestFile)
		if err != nil {
			return err
		}
		p, err := fetchURL(defaultClient(), method, u, body, headers)
		if err == nil {
			pages = append(pages, p)
		}
	}

	// 2) Offline dir
	if opts.Dir != "" {
		filepath.Walk(opts.Dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			pages = append(pages, Page{URL: path, Body: string(data)})
			return nil
		})
	}

	// 3) Input - single URL or file of URLs
	var inputTargets []string
	if opts.Input != "" {
		if isURL(opts.Input) {
			inputTargets = append(inputTargets, opts.Input)
		} else {
			lines, err := readLines(opts.Input)
			if err != nil {
				return err
			}
			for _, l := range lines {
				if isURL(l) {
					inputTargets = append(inputTargets, l)
				}
			}
		}
	}

	// 4) robots + sitemap discovery (if requested)
	var sitemapURLs []string
	if opts.Robots {
		for _, t := range inputTargets {
			if s, err := discoverSitemapsFromRobots(t); err == nil {
				sitemapURLs = append(sitemapURLs, s...)
			}
		}
		// fetch sitemap urls
		for _, s := range sitemapURLs {
			if urls, err := fetchSitemapURLs(s); err == nil {
				for _, u := range urls {
					inputTargets = append(inputTargets, u)
				}
			}
		}
	}

	// 5) Crawl/render/fetch targets
	for _, tgt := range inputTargets {
		if opts.Headless {
			// render via chromedp with 20s timeout
			if pg, err := renderPageChromedp(tgt, 20*time.Second); err == nil {
				pages = append(pages, pg)
			} else {
				// fallback to HTTP GET if chromedp fails
				if p, err := fetchURL(defaultClient(), "GET", tgt, "", nil); err == nil {
					pages = append(pages, p)
				}
			}
		} else {
			if p, err := fetchURL(defaultClient(), "GET", tgt, "", nil); err == nil {
				pages = append(pages, p)
			}
		}

		if opts.Crawl {
			crawled := crawlCollect(opts, tgt)
			pages = append(pages, crawled...)
		}
	}

	// 6) Extraction
	var allParams []string
	for _, page := range pages {
		allParams = append(allParams, extractParamsFromURL(page.URL)...)
		allParams = append(allParams, extractParamsFromHTML(page)...)
	}

	allParams = filterParams(allParams, opts.MinLength, opts.MaxLength)

	// 7) Guessing
	if opts.Guess {
		guesses := generateGuessesFromPages(pages, opts.MinLength, opts.MaxLength)
		// append to results
		allParams = append(allParams, guesses...)
		// optionally test guesses
		if opts.GuessTest && len(inputTargets) > 0 {
			// test on first input target (configurable later)
			res, _ := testGuessList(inputTargets[0], guesses, 10*time.Second)
			// dump simple report into a file next to output
			reportPath := opts.OutputFile + ".guess-report.txt"
			f, _ := os.Create(reportPath)
			defer f.Close()
			fmt.Fprintln(f, "param\tresult")
			for k, v := range res {
				fmt.Fprintf(f, "%s\t%d\n", k, v)
			}
			if !opts.Silent {
				fmt.Println("[+] Guess report saved to", reportPath)
			}
		}
	}

	allParams = filterParams(allParams, opts.MinLength, opts.MaxLength)

	if len(allParams) == 0 {
		if !opts.Silent {
			fmt.Println("[*] No parameters found.")
		}
		return nil
	}

	// Save results
	if err := writeParamsToFile(opts.OutputFile, allParams); err != nil {
		return err
	}
	if !opts.Silent {
		fmt.Printf("[+] Found %d unique parameters, saved to %s\n", len(allParams), opts.OutputFile)
	}
	return nil
}
