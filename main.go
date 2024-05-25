package main

import (
	"fmt"
	"time"

	"github.com/jassuwu/scrape-psgtech/indexer"
)

var (
	PSGTECH_JSON        = "data/psgtech.json"
	INVERTED_INDEX_JSON = "data/inverted_index.json"
	K1                  = 1.2
	B                   = 0.75
)

func main() {
	fmt.Println("-- helloworld from scrape-psgtech --")

	// scraper.Scrape()

	startIndexing := time.Now()
	err := indexer.IndexDocuments(PSGTECH_JSON, INVERTED_INDEX_JSON, K1, B)
	if err != nil {
		fmt.Println("Error indexing the documents", err)
	} else {
		fmt.Println("Indexing completed successfully in: ", time.Since(startIndexing))
	}
}
