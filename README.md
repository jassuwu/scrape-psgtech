# scraper and server for psgoogle

This repository contains two packages that performs two parts of the whole backend.

## Scraper
This will access and read the web pages and get the text after which the preprocessing occurs.

1. Crawl
2. Tokenize
3. Stopword Removal
4. Normalize
5. Stem/Lemmatize
6. Store

This will be storing 2 important information. The processed text documents, and an edgelist of all the links in some good format.

## Server
This will be public facing, initially read the information into memory (since, it's not a huge corpus in this case, this shouldn't be a huge problem), then
construct the web graph for pageranking, and store them in memory as well.

The server will take queries for search with keywords and search the corpus for matching documents and return the documents as well.
Some things that the server should be able to do are:

1. Calc and store PageRank scores
2. Process queries (maybe even some pseudo relevance feedback ?)
3. Serve requests really fast.
4. Able to receive User feedback and adjust some kind of multiplier for every document.


# Research

1. Go pkgs for all the steps above and below
2. Stemming vs. Lemmatizing
3. TFIDF vs. BM25 (vs. Self-hosted elastic-search)
4. In-memory vs. DB for storage
5. Query processing steps
