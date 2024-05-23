package scraper

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

var STARTING_LINK = "http://www.psgtech.edu"
var pageText strings.Builder
var visitedURLs = make(map[string]bool)

func Scrape() {
  c := NewCustomCollyCollector()

  c.OnRequest(onRequest)
  c.OnError(onError)
  c.OnHTML("a[href]", onAnchorTag)
  c.OnHTML("p, h1, h2, h3, h4, h5, h6, li, a, div", onTextTags)
  c.OnScraped(onScraped)

  c.Visit(STARTING_LINK)
}

func NewCustomCollyCollector() *colly.Collector {
  c := colly.NewCollector(
    colly.AllowedDomains("psgtech.edu", "www.psgtech.edu"),
  )

  /*
    * Error 1: Get "https://www.psgtech.edu/": remote error: tls: handshake failure
    * Fix for Error 1: Include the tls.TLS_RSA_WITH_RC4_128_SHA CipherSuite to the tls.Config

    * Error 2: Get "https://www.psgtech.edu/": tls: failed to verify certificate: x509: certificate signed by unknown authority
    * Fix for Error 2: Include the InsecureSkipVerify: true field in the tls.Config
  */
  c.WithTransport(
    &http.Transport{
      TLSClientConfig: &tls.Config{
        InsecureSkipVerify: true,
        CipherSuites: []uint16{
          tls.TLS_RSA_WITH_RC4_128_SHA,
        },
      },
    })

  return c
}

func onRequest (r *colly.Request) {
  fmt.Println("ABOUT TO VISIT", r.URL.String())
}

func onError (r *colly.Response, err error) {
  fmt.Println("REQUEST URL:", r.Request.URL, "FAILED WITH RESPONSE", r, "\nERROR:", err)
}

func onScraped(r *colly.Response) {
  fmt.Println("PAGE TEXT: ", pageText.String())
  pageText.Reset()
  fmt.Println("SCRAPED", r.Request.URL.String())
}

func onAnchorTag(e *colly.HTMLElement) {
  link := e.Request.AbsoluteURL(e.Attr("href"))
  pageURL := strings.Replace(link, "https://www.", "https://", 0)
  pageURL = strings.Replace(pageURL, "/index.html", "", 0)

  if pageURL != "" && !visitedURLs[pageURL] {
    visitedURLs[pageURL] = true
    e.Request.Visit(pageURL)
  }
}

func onTextTags(e *colly.HTMLElement) {
  text := strings.Join(
    strings.Fields(
      regexp.MustCompile(`[^\w\s]`).ReplaceAllString(
        strings.ToLower(e.Text),
        ""),
      ),
    " ") + " "

  if text != "" {
    pageText.WriteString(text)
  }
}
