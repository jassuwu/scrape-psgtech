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
    // colly.AllowedDomains("psgtech.edu", "www.psgtech.edu"),
  )

  c.WithTransport(
    &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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
    fmt.Println("We're getting here no?")
    link := e.Attr("href")
    fmt.Println("Link found ", e.Text, "->", link)
    e.Request.Visit(e.Request.AbsoluteURL(link))
  })


  c.Visit("http://www.psgtech.edu")

  fmt.Println("-- goodbyeworld from the scraper --")
}
