package indexer

import (
	"math"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/jassuwu/scrape-psgtech/scraper"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type BM25Score struct {
	DocumentURL string  `json:"documentURL"`
	Score       float64 `json:"score"`
}

type InvertedIndex struct {
	BM25Scores    map[string][]BM25Score
	AvgDocLength  float64
	K1            float64
	B             float64
	DocumentCount int
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

func loadScrapedDocuments(fileName string) ([]scraper.PageDocument, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	documents := []scraper.PageDocument{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&documents)
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func calculateBM25Scores(docs []scraper.PageDocument, k1, b float64) (*InvertedIndex, error) {
	termFrequency := make(map[string]map[string]int)
	documentFrequency := make(map[string]int)
	documentLength := make(map[string]int)

	totalLength := 0

	for _, doc := range docs {
		terms := strings.Fields(doc.ProcessedText)
		docLength := len(terms)
		documentLength[doc.Url] = docLength
		totalLength += docLength

		termCount := make(map[string]int)
		for _, term := range terms {
			termCount[term]++
		}

		for term, count := range termCount {
			if termFrequency[term] == nil {
				termFrequency[term] = make(map[string]int)
			}
			termFrequency[term][doc.Url] = count
			documentFrequency[term]++
		}
	}

	avgDocLen := float64(totalLength) / float64(len(docs))
	bm25Scores := make(map[string][]BM25Score)

	for term, docFreqMap := range termFrequency {
		idf := math.Log(
			(float64(len(docs)) - float64(documentFrequency[term]) + 0.5) / (float64(documentFrequency[term]) + 0.5),
		)
		for docURL, tf := range docFreqMap {
			docLen := float64(documentLength[docURL])
			score := idf * (float64(tf) * (k1 + 1)) / (float64(tf) + k1*(1-b+b*(docLen/avgDocLen)))
			bm25Scores[term] = append(
				bm25Scores[term],
				BM25Score{DocumentURL: docURL, Score: score},
			)
		}
	}

	return &InvertedIndex{
		BM25Scores:    bm25Scores,
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
	encoder.SetIndent("", "  ")
	if err != nil {
		return err
	}

	return nil
}
