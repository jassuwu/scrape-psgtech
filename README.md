# scraper, indexer and server for psgoogle

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

## Server
This will be public facing, initially read the information into memory (since, it's not a huge corpus in this case, this shouldn't be a huge problem), then
construct the web graph for pageranking, and store them in memory as well.

The server will take queries for search with keywords and search the corpus for matching documents and return the documents as well.
Some things that the server should be able to do are:

- [ ] Calc and store PageRank scores
- [x] Process queries
- [x] Serve requests really fast.
    - [x] Caching search results to serve requests faster.
- [ ] Able to receive User feedback and adjust some kind of multiplier for every document.

## Speeds
- ![Scraping and Indexing](/repoassets/scrape_and_index_speed.png "scrape_and_index_speed")
- ![Difference due to caching](/repoassets/caching_diff.png "caching_diff")

# Research

- Incorporation of PageRank scores.
- Approriate method of user relevance feedback.
