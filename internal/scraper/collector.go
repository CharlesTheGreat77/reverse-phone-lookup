package scraper

import (
	"net/http"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/imroc/req/v3"
)

// function to create a collector that impersonates a chrome browser
func createCollector() *colly.Collector {
	fakeChrome := req.DefaultClient().ImpersonateChrome()
	c := colly.NewCollector(
		colly.UserAgent(fakeChrome.Headers.Get("user-agent")),
	   )
	c.SetClient(&http.Client{Transport: fakeChrome.Transport,})
	return c
}

// set behavior of colly for each request and error response
func setCollyBehavior(c *colly.Collector) {
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("\r[*] Requesting: %s\n", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
			fmt.Printf("\n[-] Error occurred\n -> Response Code: %v\n --> Error: %v\n\n", r.StatusCode, err)
	})
}