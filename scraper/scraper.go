package scraper

import (
  "fmt"
  "net/http"
  "crypto/tls"

  "github.com/gocolly/colly/v2"
)

func Scrape() {
  fmt.Println("-- helloworld from the scraper --")
  c := colly.NewCollector(
    colly.AllowedDomains("psgtech.edu", "www.psgtech.edu"),
    )

  /*
    Error 1:  Get "https://www.psgtech.edu/": remote error: tls: handshake failure
    Fix for Error 1: Include the tls.TLS_RSA_WITH_RC4_128_SHA CipherSuite to the tls.Config

    Error 2:Something went wrong:  Get "https://www.psgtech.edu/": tls: failed to verify certificate: x509: certificate signed by unknown authority
    Fix for Error 2: Include the InsecureSkipVerify: true field in the tls.Config
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

  c.OnRequest(func(r *colly.Request) {
    fmt.Println("Visiting", r.URL.String())
  })

  c.OnError(func(_ *colly.Response, err error) {
    fmt.Println("Something went wrong: ", err)
  })

  c.OnResponse(func(r *colly.Response) {
    fmt.Println("Visited", r.Request.URL.String())
  })

  c.OnHTML("a[href]", func (e *colly.HTMLElement) {
    link := e.Attr("href")
    fmt.Println("Link found ", e.Text, "->", link)
    e.Request.Visit(e.Request.AbsoluteURL(link))
  })

  c.Visit("http://www.psgtech.edu")

  fmt.Println("-- goodbyeworld from the scraper --")
}
