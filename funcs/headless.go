package funcs

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

// renderPageChromedp visits `target` and returns the fully rendered HTML (timeouted).
// Requires Chrome/Chromium installed. If Headless false you keep using fetchURL.
func renderPageChromedp(target string, timeout time.Duration) (Page, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// create context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	cctx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	var html string
	tasks := chromedp.Tasks{
		chromedp.Navigate(target),
		// optional wait: network idle-ish
		chromedp.Sleep(500 * time.Millisecond),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	}
	if err := chromedp.Run(cctx, tasks); err != nil {
		return Page{}, err
	}
	return Page{URL: target, Body: html}, nil
}
