package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/jassuwu/scrape-psgtech/indexer"
	"github.com/jassuwu/scrape-psgtech/scraper"
	"github.com/jassuwu/scrape-psgtech/server"
)

var (
	PSGTECH_JSON        = "data/psgtech.json"
	INVERTED_INDEX_JSON = "data/inverted_index.json"
	K1                  = 1.2
	B                   = 0.75
)

func main() {
	log.Println("PROGRAM INITIATED")

	if _, err := os.Stat(PSGTECH_JSON); errors.Is(err, os.ErrNotExist) {
		startScraping := time.Now()
		scraper.Scrape()
		log.Println("Scraping completed successfully in: ", time.Since(startScraping))
	} else {
		log.Println("PSGTECH_JSON already exists. Scraping did not initiate.")
	}

	if _, err := os.Stat(INVERTED_INDEX_JSON); errors.Is(err, os.ErrNotExist) {
		startIndexing := time.Now()
		err := indexer.IndexDocuments(PSGTECH_JSON, INVERTED_INDEX_JSON, K1, B)
		if err != nil {
			log.Println("Error indexing the documents", err)
		} else {
			log.Println("Indexing completed successfully in: ", time.Since(startIndexing))
		}
	} else {
		log.Println("INVERTED_INDEX_JSON already exists. Indexing did not initiate.")
	}

	// Start the server in the end
	server.Serve()
}
