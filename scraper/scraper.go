package scraper

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

var STARTING_LINK = "http://www.psgtech.edu"
var visitedURLs = make(map[string]bool)
var punctuationRegex = regexp.MustCompile(`[^\w\s]`)

func Scrape() {
  c := NewCustomCollyCollector()

  c.OnRequest(onRequest)
  c.OnError(onError)
  c.OnResponse(onResponse)
  c.OnHTML("a[href]", onHTML)

  c.Visit(STARTING_LINK)
}

// func initFile(fName string) {
//   file, err := os.Create(fName)
//   if err != nil {
//     log.Fatalf("Cannot create file %q: %s\n", fName, err)
//   }
//   defer file.Close()
//
//   writer := csv.NewWriter(file)
//   defer writer.Flush()
// }

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
  fmt.Println("VISITED", r.URL.String())
}

func onError (r *colly.Response, err error) {
  fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
}

func onResponse(r *colly.Response) {
  // fmt.Println("Visited", r.Request.URL.String())
}

func onHTML(e *colly.HTMLElement) {
  // Page URL Cleaning and Visited URL Check
  pageURL := strings.Replace(e.Request.URL.String(), "https://www.", "https://", 1)
  if visitedURLs[pageURL] {
    return
  }
  visitedURLs[pageURL] = true

  // Page Title Extraction
  pageTitle := e.DOM.Find("title").Text()

  // Page Text Extraction and Punctuation Removal
  var pageText strings.Builder
  e.ForEach("p, h1, h2, h3, h4, h5, h6, li, a, div", func(_ int, el *colly.HTMLElement) {
    pageText.WriteString(el.Text)
    pageText.WriteString(" ")
  })

  cleanText := punctuationRegex.ReplaceAllString(pageText.String(), "")

  // URL Parsing and Processing
  parsedURL, _ := url.Parse(pageURL)
  parts := strings.Split(parsedURL.Host+parsedURL.Path, "/")
  var processedURL []string
  for _, part := range parts {
    if part != "" {
      processedURL = append(processedURL, part)
    }
  }

  // Processing Text and Title
  // In this example, we're not handling stopwords and stemming
  // You can use external libraries like github.com/kljensen/snowball for stemming
  processedText := strings.Fields(strings.ToLower(cleanText))
  processedTitle := strings.Fields(strings.ToLower(pageTitle))

  // Extracting Links
  var links []string
  e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
    link := el.Request.AbsoluteURL(el.Attr("href"))
    cleanedLink := strings.Replace(link, "https://www.", "https://", 1)
    links = append(links, cleanedLink)
  })

  // Yielding Results
  result := map[string]any{
    "url":   pageURL,
    "title": processedTitle,
    "text":  processedText,
    "links": links,
  }

  fmt.Println(result)
  e.Request.Visit(pageURL)
}
