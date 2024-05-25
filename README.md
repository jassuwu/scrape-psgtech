# scraper and server for psgoogle

This repository contains packages that performs parts of the whole backend.

## Scraper
This will access and read the web pages and get the text after which the preprocessing occurs.

- [x] Crawl - `github.com/gocolly/colly/v2`
- [x] Tokenize
- [x] Stopword Removal
- [x] Normalize
- [x] Stem/Lemmatize
- [x] Store as JSON

This will be storing 2 important information. The processed text documents, and an edgelist of all the links in some good format.

## Indexer
This will take the data stores in the JSON and process it as per the chosen IR algorithm and store it in the appropriate data structure, either in-memory or DB.

- [x] Choose the IR algorithm & data structure. (Chosen Okapi BM25)
- [x] Perform the algorithm and store in the correct format. (Stored in JSON as an inverted index.)

Testrun of the first iteration of the indexer and scraper:
![Testrun](/assets/testrun.png "testrun")

## Server
This will be public facing, initially read the information into memory (since, it's not a huge corpus in this case, this shouldn't be a huge problem), then
construct the web graph for pageranking, and store them in memory as well.

The server will take queries for search with keywords and search the corpus for matching documents and return the documents as well.
Some things that the server should be able to do are:

- [ ] Calc and store PageRank scores
- [ ] Process queries (maybe even some pseudo relevance feedback ?)
- [ ] Serve requests really fast.
- [ ] Able to receive User feedback and adjust some kind of multiplier for every document.

# Research

5. Query processing steps
