package funcs

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Page struct {
	URL  string
	Body string
}

func buildHTTPClient(proxyStr string) (*http.Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // bug-bounty style, ignore TLS issues
		Proxy:           nil,
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
	}

	if proxyStr != "" {
		pURL, err := url.Parse(proxyStr)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy: %w", err)
		}
		transport.Proxy = http.ProxyURL(pURL)
	}

	client := &http.Client{
		Timeout:   20 * time.Second,
		Transport: transport,
	}
	return client, nil
}

func fetchURL(client *http.Client, method, rawURL, body string, headers http.Header) (Page, error) {
	reqBody := io.Reader(nil)
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, rawURL, reqBody)
	if err != nil {
		return Page{}, err
	}

	for k, vals := range headers {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return Page{}, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Page{}, err
	}

	return Page{URL: rawURL, Body: string(b)}, nil
}

func parseRawRequestFile(path string) (method, urlStr string, headers http.Header, body string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", nil, "", err
	}
	reader := bufio.NewReader(bytes.NewReader(data))
	startLine, err := reader.ReadString('\n')
	if err != nil {
		return "", "", nil, "", fmt.Errorf("invalid request file: %w", err)
	}
	parts := strings.Fields(startLine)
	if len(parts) < 2 {
		return "", "", nil, "", fmt.Errorf("invalid request start line: %q", startLine)
	}
	method = parts[0]
	target := parts[1]

	headers = make(http.Header)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", "", nil, "", err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		hp := strings.SplitN(line, ":", 2)
		if len(hp) != 2 {
			continue
		}
		k := strings.TrimSpace(hp[0])
		v := strings.TrimSpace(hp[1])
		headers.Add(k, v)
	}

	bodyBuf, err := io.ReadAll(reader)
	if err != nil {
		return "", "", nil, "", err
	}
	body = string(bodyBuf)

	// Try to reconstruct full URL from Host header + path
	if !isURL(target) {
		host := headers.Get("Host")
		if host != "" {
			target = "https://" + host + target
		}
	}

	return method, target, headers, body, nil
}

func defaultClient() *http.Client {
	c, err := buildHTTPClient("")
	if err != nil {
		// fallback to a basic client
		return &http.Client{}
	}
	return c
}
