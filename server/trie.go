package server

import (
	"log"
	"time"

	"github.com/jassuwu/scrape-psgtech/indexer"
)

type TrieNode struct {
	Children map[rune]*TrieNode
	IsWord   bool
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		Children: make(map[rune]*TrieNode),
		IsWord:   false,
	}
}

type Trie struct {
	Root *TrieNode
}

func newTrie() *Trie {
	return &Trie{
		Root: newTrieNode(),
	}
}

func (t *Trie) insert(word string) {
	curr := t.Root
	for _, char := range word {
		if _, hit := curr.Children[char]; !hit {
			curr.Children[char] = newTrieNode()
		}
		curr = curr.Children[char]
	}
	curr.IsWord = true
}

func (t *Trie) findWordsWithPrefix(prefix string) []string {
	if prefix == "" {
		return []string{}
	}
	curr := t.Root
	for _, char := range prefix {
		if _, hit := curr.Children[char]; !hit {
			return []string{}
		}
		curr = curr.Children[char]
	}
	return collectWords(curr, prefix)
}

func collectWords(node *TrieNode, prefix string) []string {
	var words []string
	if node.IsWord {
		words = append(words, prefix)
	}
	for char, child := range node.Children {
		words = append(words, collectWords(child, prefix+string(char))...)
	}
	return words
}

func makeTrie(idx *indexer.InvertedIndex) *Trie {
	trie := newTrie()
	start := time.Now()
	for word := range idx.BM25Scores {
		trie.insert(word)
	}
	log.Println("Autocompletion Trie Initialized in: ", time.Since(start))
	return trie
}
