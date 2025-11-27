package funcs

import (
	"net/url"
	"regexp"
	"sort"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

var (
	attrURLSelector = []string{"a[href]", "link[href]", "script[src]", "img[src]", "iframe[src]"}
	formSelector    = "form input[name], form select[name], form textarea[name]"
	nameAttrRegex   = regexp.MustCompile(`\bname\s*=\s*["']?([a-zA-Z0-9_\-\.]+)["']?`)
)

func extractParamsFromURL(raw string) []string {
	u, err := url.Parse(raw)
	if err != nil {
		return nil
	}
	var names []string
	for name := range u.Query() {
		names = append(names, name)
	}
	return names
}

func extractParamsFromHTML(page Page) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page.Body))
	if err != nil {
		return nil
	}
	params := make(map[string]struct{})

	// 1) links / src URLs
	for _, sel := range attrURLSelector {
		doc.Find(sel).Each(func(_ int, s *goquery.Selection) {
			for _, attr := range []string{"href", "src"} {
				if v, ok := s.Attr(attr); ok {
					for _, p := range extractParamsFromURL(v) {
						params[p] = struct{}{}
					}
				}
			}
		})
	}

	// 2) form field names
	doc.Find(formSelector).Each(func(_ int, s *goquery.Selection) {
		if name, ok := s.Attr("name"); ok {
			name = strings.TrimSpace(name)
			if name != "" {
				params[name] = struct{}{}
			}
		}
	})

	// 3) fallback regex scan for name="..."
	matches := nameAttrRegex.FindAllStringSubmatch(page.Body, -1)
	for _, m := range matches {
		if len(m) > 1 {
			params[m[1]] = struct{}{}
		}
	}

	var list []string
	for p := range params {
		list = append(list, p)
	}
	sort.Strings(list)
	return list
}

func filterParams(params []string, minLen, maxLen int) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, p := range params {
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		l := len(p)
		if minLen > 0 && l < minLen {
			continue
		}
		if maxLen > 0 && l > maxLen {
			continue
		}
		out = append(out, p)
	}
	return out
}
