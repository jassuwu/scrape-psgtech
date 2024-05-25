package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jassuwu/scrape-psgtech/indexer"
	"github.com/jassuwu/scrape-psgtech/scraper"
)

var (
	PSGTECH_JSON        = "data/psgtech.json"
	INVERTED_INDEX_JSON = "data/inverted_index.json"
	K1                  = 1.2
	B                   = 0.75
)

func main() {
	startProgram := time.Now()
	fmt.Println("SCRAPING INITIATED")

	startScraping := time.Now()
	scraper.Scrape()
	log.Println("Scraping completed successfully in: ", time.Since(startScraping))

	startIndexing := time.Now()
	err := indexer.IndexDocuments(PSGTECH_JSON, INVERTED_INDEX_JSON, K1, B)
	if err != nil {
		log.Println("Error indexing the documents", err)
	} else {
		log.Println("Indexing completed successfully in: ", time.Since(startIndexing))
	}

	log.Println("Program completed successfully in: ", time.Since(startProgram))
}
