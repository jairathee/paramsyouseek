package funcs

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// HeaderList implements flag.Value for repeated -H flags
type HeaderList struct {
	Headers http.Header
}

func (h *HeaderList) String() string {
	if h == nil || h.Headers == nil {
		return ""
	}
	var b strings.Builder
	first := true
	for k, vals := range h.Headers {
		for _, v := range vals {
			if !first {
				b.WriteString(", ")
			}
			first = false
			b.WriteString(fmt.Sprintf("%s: %s", k, v))
		}
	}
	return b.String()
}

func (h *HeaderList) Set(value string) error {
	if h.Headers == nil {
		h.Headers = make(http.Header)
	}
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid header, expected Name: Value, got %q", value)
	}
	name := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	if name == "" {
		return errors.New("header name cannot be empty")
	}
	h.Headers.Add(name, val)
	return nil
}

type Options struct {
	Input            string
	Dir              string
	RequestFile      string
	Threads          int
	DelaySeconds     int
	Crawl            bool
	Depth            int
	CrawlDurationStr string
	CrawlDuration    time.Duration
	Headless         bool // render JS
	HeaderFlags      HeaderList
	Method           string
	Body             string
	Proxy            string

	OutputFile string
	MaxLength  int
	MinLength  int
	Silent     bool

	DisableUpdateCheck bool

	// New flags
	Guess          bool     // generate guesses from pages
	GuessTest      bool     // actively test guesses against server
	AllowedDomains []string // optional domain whitelist: if set, only crawl domains in this list
	Robots         bool     // obey robots.txt and discover sitemaps
}

func (o *Options) Normalize() error {
	if o.Threads <= 0 {
		o.Threads = 1
	}
	if o.MaxLength <= 0 {
		o.MaxLength = 30
	}
	if o.CrawlDurationStr != "" {
		d, err := time.ParseDuration(o.CrawlDurationStr)
		if err != nil {
			return fmt.Errorf("invalid crawl duration %q: %w", o.CrawlDurationStr, err)
		}
		o.CrawlDuration = d
	}
	// Normalize method
	if o.Method == "" {
		o.Method = "GET"
	}
	o.Method = strings.ToUpper(o.Method)

	// Normalize AllowedDomains: lowercase and trim.
	for i, d := range o.AllowedDomains {
		o.AllowedDomains[i] = strings.ToLower(strings.TrimSpace(d))
	}
	return nil
}
