package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jairathee/paramsyouseek/funcs"
)

func main() {
	opts := funcs.Options{}

	flag.StringVar(&opts.Input, "u", "", "Input [Filename | URL]")
	flag.StringVar(&opts.Input, "url", "", "Input [Filename | URL] (alias)")

	flag.StringVar(&opts.Dir, "dir", "", "Stored requests/responses directory (offline)")
	flag.StringVar(&opts.Dir, "directory", "", "Stored requests/responses directory (offline) (alias)")

	flag.StringVar(&opts.RequestFile, "r", "", "Raw HTTP request file")
	flag.StringVar(&opts.RequestFile, "request", "", "Raw HTTP request file (alias)")

	flag.IntVar(&opts.Threads, "t", 1, "Number of threads")
	flag.IntVar(&opts.Threads, "thread", 1, "Number of threads (alias)")

	flag.IntVar(&opts.DelaySeconds, "rd", 0, "Request delay between each request in seconds")
	flag.IntVar(&opts.DelaySeconds, "delay", 0, "Request delay between each request in seconds (alias)")

	flag.BoolVar(&opts.Crawl, "c", false, "Crawl pages to extract their parameters")
	flag.BoolVar(&opts.Crawl, "crawl", false, "Crawl pages to extract their parameters (alias)")

	flag.IntVar(&opts.Depth, "d", 2, "Maximum crawl depth")
	flag.IntVar(&opts.Depth, "depth", 2, "Maximum crawl depth (alias)")

	flag.StringVar(&opts.CrawlDurationStr, "ct", "", "Maximum crawl duration (e.g. 30s, 2m)")
	flag.StringVar(&opts.CrawlDurationStr, "crawl-duration", "", "Maximum crawl duration (alias)")

	flag.BoolVar(&opts.Headless, "hl", false, "Headless mode (render JS; basic support)")
	flag.BoolVar(&opts.Headless, "headless", false, "Headless mode (alias)")

	flag.BoolVar(&opts.Guess, "guess", false, "Generate guessed parameters")
	flag.BoolVar(&opts.GuessTest, "guesstest", false, "Test guessed parameters for reflections")

	flag.BoolVar(&opts.Robots, "robots", false, "Parse robots.txt and sitemap.xml for URLs")

	flag.Var(&opts.HeaderFlags, "H", "Header \"Name: Value\" (can be repeated)")
	flag.Var(&opts.HeaderFlags, "header", "Header \"Name: Value\" (alias)")

	flag.StringVar(&opts.Method, "X", "GET", "HTTP method")
	flag.StringVar(&opts.Method, "method", "GET", "HTTP method (alias)")

	flag.StringVar(&opts.Body, "b", "", "POST data")
	flag.StringVar(&opts.Body, "body", "", "POST data (alias)")

	flag.StringVar(&opts.Proxy, "x", "", "Proxy URL (SOCKS5 or HTTP)")
	flag.StringVar(&opts.Proxy, "proxy", "", "Proxy URL (alias)")

	flag.StringVar(&opts.OutputFile, "o", "parameters.txt", "Output file")
	flag.StringVar(&opts.OutputFile, "output", "parameters.txt", "Output file (alias)")

	flag.IntVar(&opts.MaxLength, "xl", 30, "Maximum length of parameter names")
	flag.IntVar(&opts.MaxLength, "max-length", 30, "Maximum length of parameter names (alias)")

	flag.IntVar(&opts.MinLength, "nl", 0, "Minimum length of parameter names")
	flag.IntVar(&opts.MinLength, "min-length", 0, "Minimum length of parameter names (alias)")

	flag.BoolVar(&opts.Silent, "silent", false, "Disable banner and extra logs")

	flag.BoolVar(&opts.DisableUpdateCheck, "duc", false, "Disable update check (no-op placeholder)")
	flag.BoolVar(&opts.DisableUpdateCheck, "disable-update-check", false, "Disable update check (alias)")

	flag.Parse()

	if opts.Input == "" && opts.Dir == "" && opts.RequestFile == "" {
		fmt.Fprintln(os.Stderr, "[-] You must provide -u, -dir or -r")
		flag.Usage()
		os.Exit(1)
	}

	if err := opts.Normalize(); err != nil {
		fmt.Fprintf(os.Stderr, "[-] options error: %v\n", err)
		os.Exit(1)
	}

	if err := funcs.Run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "[-] error: %v\n", err)
		os.Exit(1)
	}
}
