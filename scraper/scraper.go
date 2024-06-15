package scraper

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

var (
	PSGTECH_JSON_FILE_PATH = "data/psgtech.json"
	STARTING_LINK          = "http://www.psgtech.edu/"
)

var (
	pageTitle   string
	pageText    strings.Builder
	pageLinks   []string
	visitedURLs = make(map[string]bool)
)

type PageDocument struct {
	Url           string   `json:"url"`
	Title         string   `json:"title"`
	ProcessedText string   `json:"processedText"`
	Links         []string `json:"links"`
}

var pageDocuments []PageDocument

func Scrape() {
	c := NewCustomCollyCollector()

	// c.OnRequest(onRequest)
	c.OnError(onError)
	c.OnHTML("a[href]", onAnchorTag)
	c.OnHTML("title", onTitleTag)
	c.OnHTML("p, h1, h2, h3, h4, h5, h6, li, a, div", onTextTags)
	c.OnScraped(onScraped)

	initJSONAndStartVisiting(STARTING_LINK, c.Visit)
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

func onRequest(r *colly.Request) {
	fmt.Println("ABOUT TO VISIT", r.URL.String())
}

func onError(r *colly.Response, err error) {
	log.Println("ERROR:", err, "IN URL", r.Request.URL, "FAILED WITH RESPONSE", r.StatusCode)
}

func onScraped(r *colly.Response) {
	moreText := " " + strings.Join(
		strings.Fields(
			regexp.MustCompile(
				`[^\w\s]`,
			).ReplaceAllString(strings.ToLower(r.Request.URL.String()+" "+pageTitle), " "),
		),
		" ")
	pageText.WriteString(StemText(RemoveStopWords(moreText)))

	pageDocument := PageDocument{
		Url:           r.Request.URL.String(),
		Title:         pageTitle,
		ProcessedText: strings.Trim(pageText.String(), " "),
		Links:         pageLinks,
	}

	pageText.Reset()
	pageLinks = nil
	log.Println("SCRAPED", r.Request.URL.String())
	appendToJSON(pageDocument)
}

func onTitleTag(e *colly.HTMLElement) {
	pageTitle = e.Text
}

func onAnchorTag(e *colly.HTMLElement) {
	link := e.Request.AbsoluteURL(e.Attr("href"))
	pageURL := strings.Replace(link, "https://www.", "https://", 1)
	pageURL = strings.Replace(pageURL, "index.html", "", 1)

	shouldVisit := pageURL != "" &&
		strings.HasPrefix(pageURL, "http") &&
		!visitedURLs[pageURL] &&
		!hasExcludedExtension(pageURL)

	if shouldVisit {
		visitedURLs[pageURL] = true
		pageLinks = append(pageLinks, pageURL)
		e.Request.Visit(pageURL)
	}
}

func onTextTags(e *colly.HTMLElement) {
	text := StemText(
		RemoveStopWords(
			NormalizeText(
				e.Text,
			),
		),
	)

	if text != "" {
		pageText.WriteString(text)
	}
}
