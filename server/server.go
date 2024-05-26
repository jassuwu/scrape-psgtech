package server

import (
	"fmt"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var (
	json                                  = jsoniter.ConfigCompatibleWithStandardLibrary
	documents, invertedIndex, dataLoadErr = loadData(
		"data/psgtech.json",
		"data/inverted_index.json",
	)
)

func Serve() {
	if dataLoadErr != nil {
		log.Fatal(dataLoadErr)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", root)
	mux.HandleFunc("GET /ping", ping)
	mux.HandleFunc("GET /health", ping)
	mux.HandleFunc("GET /search", search)

	http.ListenAndServe(":8000", mux)
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

	q := r.FormValue("q")
	rankedDocuments := rankDocuments(q, documents, invertedIndex)

	data, err := json.MarshalToString(
		map[string]any{"results": rankedDocuments, "error": false},
	)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintln(w, data)
}
