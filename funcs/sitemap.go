package funcs

import (
	"encoding/xml"
	"io"
	"strings"
)

type urlset struct {
	Locations        []string `xml:"url>loc"`
	SitemapLocations []string `xml:"sitemap>loc"`
}

// discoverSitemapsFromRobots fetches robots.txt and returns discovered sitemap URLs.
func discoverSitemapsFromRobots(target string) ([]string, error) {
	client := defaultClient()
	robots := joinURLPath(target, "/robots.txt")
	resp, err := client.Get(robots)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var sitemaps []string
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "sitemap:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				sitemaps = append(sitemaps, strings.TrimSpace(parts[1]))
			}
		}
	}
	return sitemaps, nil
}

func joinURLPath(base string, extra string) string {
	if strings.HasSuffix(base, "/") {
		base = strings.TrimRight(base, "/")
	}
	if strings.HasPrefix(extra, "/") {
		extra = strings.TrimLeft(extra, "/")
	}
	return base + "/" + extra
}

// fetchSitemapURLs fetches an XML sitemap and returns contained <loc> URLs.
func fetchSitemapURLs(smap string) ([]string, error) {
	client := defaultClient()
	resp, err := client.Get(smap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var us urlset
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Try to decode <urlset> or <sitemapindex>
	if err := xml.Unmarshal(body, &us); err != nil {
		// fallback crude: find <loc> tags by string
		out := []string{}
		s := string(body)
		for {
			i := strings.Index(s, "<loc>")
			if i == -1 {
				break
			}
			j := strings.Index(s[i:], "</loc>")
			if j == -1 {
				break
			}
			loc := s[i+5 : i+j]
			out = append(out, strings.TrimSpace(loc))
			s = s[i+j+6:]
		}
		return out, nil
	}
	out := append(us.Locations, us.SitemapLocations...)
	return out, nil
}
