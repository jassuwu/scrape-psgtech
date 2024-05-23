package scraper

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

var STARTING_LINK = "http://www.psgtech.edu"
var pageTitle string
var pageText strings.Builder
var pageLinks []string
var visitedURLs = make(map[string]bool)

type PageDocument struct {
  Url, Title, ProcessedText string
  Links []string
}

var pageDocuments []PageDocument

func Scrape() {
  c := NewCustomCollyCollector()

  c.OnRequest(onRequest)
  c.OnError(onError)
  c.OnHTML("a[href]", onAnchorTag)
  c.OnHTML("title", onTitleTag)
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
  moreText := strings.Join(
		strings.Fields(
			regexp.MustCompile(
				`[^\w\s]`,
			).ReplaceAllString(strings.ToLower(r.Request.URL.String() + " " + pageTitle), " "),
		),
	" ")

  pageText.WriteString(moreText)

  pageDocument := PageDocument{
    Url: r.Request.URL.String(),
    Title: pageTitle,
    ProcessedText: strings.Trim(pageText.String(), " "),
    Links: pageLinks,
  }

  pageDocuments = append(pageDocuments, pageDocument)
  pageText.Reset()
  fmt.Println("SCRAPED", r.Request.URL.String())
  saveToJSON("psgtech.json")
}

func onTitleTag(e *colly.HTMLElement) {
  pageTitle = e.Text
}

func onAnchorTag(e *colly.HTMLElement) {
  link := e.Request.AbsoluteURL(e.Attr("href"))
  pageURL := strings.Replace(link, "https://www.", "https://", 0)
  pageURL = strings.Replace(pageURL, "/index.html", "", 0)

  if pageURL != "" && strings.HasPrefix(pageURL, "http") && !visitedURLs[pageURL] {
    visitedURLs[pageURL] = true
    pageLinks = append(pageLinks, pageURL)
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

func saveToJSON(fileName string) {
  file, err := os.Create(fileName)
  if err != nil {
    fmt.Println("JSON File couldn't be created: ", err)
  }
  defer file.Close()

  encoder := json.NewEncoder(file)
  encoder.SetIndent("", "  ")
  err = encoder.Encode(pageDocuments)
  if err != nil {
    fmt.Println("Couldn't encode data to JSON: ", err)
  }
}
