package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/jassuwu/scrape-psgtech/indexer"
	"github.com/jassuwu/scrape-psgtech/scraper"
)

var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	documents     map[string]scraper.PageDocument
	invertedIndex *indexer.InvertedIndex
	dataLoadErr   error
	searchCache   *Cache
)

func Serve() {
	documents, invertedIndex, dataLoadErr = loadData(
		"data/psgtech.json",
		"data/inverted_index.json",
	)

	if dataLoadErr != nil {
		log.Fatal(dataLoadErr)
		return
	}

	searchCache = NewCache()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", root)
	mux.HandleFunc("GET /ping", ping)
	mux.HandleFunc("GET /health", ping)
	mux.HandleFunc("GET /search", search)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8000"
	}
	http.ListenAndServe(addr, mux)
}

func root(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	data, err := json.MarshalToString(
		map[string]any{"message": "welcome to psgoogle's backend.", "error": false},
	)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintln(w, data)
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	data, err := json.MarshalToString(
		map[string]any{"message": "we healthy! or should i say pong :)", "error": false},
	)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintln(w, data)
}

func search(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	start := time.Now()

	q := r.FormValue("q")
	cacheKey := "search:" + q
	cacheHit := false

	rankedDocuments := []RankedDocument{}
	if cachedResults, err := searchCache.Get(cacheKey); err == nil {
		rankedDocuments = cachedResults
		cacheHit = true
	} else {
		rankedDocuments = rankDocuments(q, documents, invertedIndex)
		if _, err := searchCache.Set(cacheKey, rankedDocuments); err != nil {
			log.Println("Error caching results for search: ", q)
		}
	}

	data, err := json.MarshalToString(
		map[string]any{"results": rankedDocuments, "error": false},
	)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	// log to stdout
	if cacheHit {
		log.Printf("/search?q=\"%s\" took %v. (from cache)\n", q, time.Since(start))
	} else {
		log.Printf("/search?q=\"%s\" took %v.\n", q, time.Since(start))
	}
	fmt.Fprintln(w, data)
}
