package scraper

import (
  "fmt"
  "strings"
  "regexp"
  "net/http"
	"crypto/tls"

	"github.com/gocolly/colly/v2"
)

var STARTING_LINK = "http://www.psgtech.edu"

func Scrape() {
  c := NewCustomCollyCollector()

  c.OnRequest(onRequest)
  c.OnError(onError)
  c.OnResponse(onResponse)
  c.OnHTML("a[href]", onAnchorTag)
  c.OnHTML("p, h1, h2, h3, h4, h5, h6, li, a, div", onTextTags)

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

func onAnchorTag(e *colly.HTMLElement) {
  pageURL := strings.Replace(e.Request.URL.String(), "https://www.", "https://", 1)
  e.Request.Visit(pageURL)
}

func onTextTags(e *colly.HTMLElement) {
  text := strings.Join(
    strings.Fields(
      regexp.MustCompile(`[^\w\s]`).ReplaceAllString(
        strings.ToLower(e.Text),
        ""),
      ),
    " ")
  if text != "" {
    fmt.Println(e.Name, "->", text)
  }
}
