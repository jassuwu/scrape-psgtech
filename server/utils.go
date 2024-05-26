package server

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jassuwu/scrape-psgtech/indexer"
	"github.com/jassuwu/scrape-psgtech/scraper"
)

type RankedDocument struct {
	Url   string  `json:"url"`
	Title string  `json:"title"`
	Score float64 `json:"score"`
}

func rankDocuments(
	query string,
	docs []scraper.PageDocument,
	idx *indexer.InvertedIndex,
) []RankedDocument {
	queryTerms := strings.Fields(processQuery(query))
	docScores := make(map[string]float64)

	for _, term := range queryTerms {
		if scores, found := idx.BM25Scores[term]; found {
			for _, score := range scores {
				docScores[score.DocumentURL] += score.Score
			}
		}
	}

	rankedDocs := make([]RankedDocument, 0, len(docScores))
	for url, score := range docScores {
		var title string
		for _, doc := range docs {
			if doc.Url == url {
				title = doc.Title
			}
			break
		}
		rankedDocs = append(rankedDocs, RankedDocument{
			Url:   url,
			Title: title,
			Score: score,
		})
	}

	// Reverse sort the ranked documents w.r.t. the BM25 score
	sort.Slice(rankedDocs, func(i, j int) bool {
		return rankedDocs[i].Score > rankedDocs[j].Score
	})

	return rankedDocs
}

func processQuery(query string) string {
	return scraper.StemText(
		scraper.RemoveStopWords(
			scraper.NormalizeText(
				query,
			),
		),
	)
}

func loadData(
	psgtechFilePath, invertedIndexFilePath string,
) ([]scraper.PageDocument, *indexer.InvertedIndex, error) {
	psgtechFile, err := os.OpenFile(psgtechFilePath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	defer psgtechFile.Close()

	documents := []scraper.PageDocument{}
	decoder := json.NewDecoder(psgtechFile)
	err = decoder.Decode(&documents)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	invertedIndexFile, err := os.OpenFile(invertedIndexFilePath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	defer invertedIndexFile.Close()

	var invIdx indexer.InvertedIndex
	decoder = json.NewDecoder(invertedIndexFile)
	err = decoder.Decode(&invIdx)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	return documents, &invIdx, nil
}
