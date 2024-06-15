package indexer

import (
	"math"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/jassuwu/scrape-psgtech/scraper"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type IndexedWord struct {
	DocumentURL   string   `json:"documentURL"`
	BM25Score     float64  `json:"bm25score"`
	OriginalWords []string `json:"originalWords"`
}

type InvertedIndex struct {
	IndexedWords  map[string][]IndexedWord `json:"indexedWords"`
	AvgDocLength  float64                  `json:"avgDocLength"`
	K1            float64                  `json:"k1"`
	B             float64                  `json:"b"`
	DocumentCount int                      `json:"documentCount"`
}

func IndexDocuments(inputF, outputF string, k1, b float64) error {
	// load the documents
	documents, err := loadScrapedDocuments(inputF)
	if err != nil {
		return err
	}

	// Calculate the BM25 scores into the Index
	idx, err := calculateBM25Scores(documents, k1, b)
	if err != nil {
		return err
	}

	// Save the inverted index with BM25 scores
	err = idx.save(outputF)
	if err != nil {
		return err
	}

	return nil
}

func loadScrapedDocuments(fileName string) (map[string]scraper.PageDocument, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	documents := make(map[string]scraper.PageDocument)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&documents)
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func calculateBM25Scores(
	docs map[string]scraper.PageDocument,
	k1, b float64,
) (*InvertedIndex, error) {
	termFrequency := make(map[string]map[string]int)
	documentFrequency := make(map[string]int)
	documentLength := make(map[string]int)

	totalLength := 0

	for url, doc := range docs {
		terms := strings.Fields(doc.ProcessedText)
		docLength := len(terms)
		documentLength[url] = docLength
		totalLength += docLength

		termCount := make(map[string]int)
		for _, term := range terms {
			termCount[term]++
		}

		for term, count := range termCount {
			if termFrequency[term] == nil {
				termFrequency[term] = make(map[string]int)
			}
			termFrequency[term][url] = count
			documentFrequency[term]++
		}
	}

	avgDocLen := float64(totalLength) / float64(len(docs))
	indexedWords := make(map[string][]IndexedWord)

	for term, docFreqMap := range termFrequency {
		idf := math.Log(
			(float64(len(docs)) - float64(documentFrequency[term]) + 0.5) / (float64(documentFrequency[term]) + 0.5),
		)
		for docURL, tf := range docFreqMap {
			docLen := float64(documentLength[docURL])
			score := idf * (float64(tf) * (k1 + 1)) / (float64(tf) + k1*(1-b+b*(docLen/avgDocLen)))
			indexedWords[term] = append(
				indexedWords[term],
				IndexedWord{DocumentURL: docURL, BM25Score: score},
			)
		}
	}

	return &InvertedIndex{
		IndexedWords:  indexedWords,
		AvgDocLength:  avgDocLen,
		K1:            k1,
		B:             b,
		DocumentCount: len(docs),
	}, nil
}

func (idx *InvertedIndex) save(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(idx)
	if err != nil {
		return err
	}

	return nil
}
