package funcs

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

var builtinSeeds = []string{
	"id", "id[]", "q", "query", "search", "s", "page", "limit",
	"user", "username", "email", "token", "auth", "callback",
	"redirect", "next", "return", "session", "itemId", "productId",
}

var jsVarRE = regexp.MustCompile(`(?m)(?:var|let|const)\s+([A-Za-z0-9_\-]+)`)
var wordRE = regexp.MustCompile(`([A-Za-z][A-Za-z0-9_]{1,30})`)

// generateGuessesFromPages returns unique, filtered guesses.
func generateGuessesFromPages(pages []Page, minLen, maxLen int) []string {
	set := map[string]struct{}{}
	// seed
	for _, s := range builtinSeeds {
		if len(s) >= minLen && (maxLen == 0 || len(s) <= maxLen) {
			set[s] = struct{}{}
		}
	}
	for _, p := range pages {
		// find JS vars
		for _, m := range jsVarRE.FindAllStringSubmatch(p.Body, -1) {
			if len(m) > 1 {
				name := m[1]
				if len(name) >= minLen && (maxLen == 0 || len(name) <= maxLen) {
					set[name] = struct{}{}
				}
			}
		}
		// find words in page text
		for _, m := range wordRE.FindAllStringSubmatch(p.Body, -1) {
			if len(m) > 1 {
				name := strings.ToLower(m[1])
				if len(name) >= minLen && (maxLen == 0 || len(name) <= maxLen) {
					set[name] = struct{}{}
				}
			}
		}
	}
	// collect, sort
	var out []string
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// write guesses to file (one per line)
func writeGuesses(path string, guesses []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, g := range guesses {
		fmt.Fprintln(w, g)
	}
	return w.Flush()
}

// testGuessList (optional) will try simple GET requests with param=value and record responses
func testGuessList(target string, guesses []string, timeout time.Duration) (map[string]int, error) {
	client := defaultClient()
	results := map[string]int{}
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	original := ""
	// fetch baseline (no param)
	if page, err := fetchURL(client, "GET", target, "", nil); err == nil {
		original = page.Body
	}
	for _, g := range guesses {
		q := u.Query()
		q.Set(g, "paramsyouseek_test")
		u.RawQuery = q.Encode()
		page, err := fetchURL(client, "GET", u.String(), "", nil)
		if err != nil {
			results[g] = -1
			continue
		}
		// simplistic difference detection: if body differs in length or contained token
		if page.Body != original && strings.Contains(page.Body, "paramsyouseek_test") {
			results[g] = 2 // reflected token
		} else if len(page.Body) != len(original) {
			results[g] = 1 // response changed
		} else {
			results[g] = 0 // no obvious change
		}
		// tiny delay to avoid flooding
		time.Sleep(200 * time.Millisecond)
	}
	return results, nil
}
