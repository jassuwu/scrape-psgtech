package scraper

import (
	"regexp"
	"strings"

	"github.com/kljensen/snowball"
)

var STOP_WORDS = map[string]struct{}{
	"i": {}, "me": {}, "my": {}, "myself": {}, "we": {}, "our": {}, "ours": {}, "ourselves": {},
	"you": {}, "you're": {}, "you've": {}, "you'll": {}, "you'd": {}, "your": {}, "yours": {},
	"yourself": {}, "yourselves": {}, "he": {}, "him": {}, "his": {}, "himself": {}, "she": {},
	"she's": {}, "her": {}, "hers": {}, "herself": {}, "it": {}, "it's": {}, "its": {}, "itself": {},
	"they": {}, "them": {}, "their": {}, "theirs": {}, "themselves": {}, "what": {}, "which": {},
	"who": {}, "whom": {}, "this": {}, "that": {}, "that'll": {}, "these": {}, "those": {},
	"am": {}, "is": {}, "are": {}, "was": {}, "were": {}, "be": {}, "been": {}, "being": {},
	"have": {}, "has": {}, "had": {}, "having": {}, "do": {}, "does": {}, "did": {}, "doing": {},
	"a": {}, "an": {}, "the": {}, "and": {}, "but": {}, "if": {}, "or": {}, "because": {},
	"as": {}, "until": {}, "while": {}, "of": {}, "at": {}, "by": {}, "for": {}, "with": {},
	"about": {}, "against": {}, "between": {}, "into": {}, "through": {}, "during": {}, "before": {},
	"after": {}, "above": {}, "below": {}, "to": {}, "from": {}, "up": {}, "down": {}, "in": {},
	"out": {}, "on": {}, "off": {}, "over": {}, "under": {}, "again": {}, "further": {}, "then": {},
	"once": {}, "here": {}, "there": {}, "when": {}, "where": {}, "why": {}, "how": {}, "all": {},
	"any": {}, "both": {}, "each": {}, "few": {}, "more": {}, "most": {}, "other": {}, "some": {},
	"such": {}, "no": {}, "nor": {}, "not": {}, "only": {}, "own": {}, "same": {}, "so": {},
	"than": {}, "too": {}, "very": {}, "s": {}, "t": {}, "can": {}, "will": {}, "just": {},
	"don": {}, "don't": {}, "should": {}, "should've": {}, "now": {}, "d": {}, "ll": {}, "m": {},
	"o": {}, "re": {}, "ve": {}, "y": {}, "ain": {}, "aren": {}, "aren't": {}, "couldn": {},
	"couldn't": {}, "didn": {}, "didn't": {}, "doesn": {}, "doesn't": {}, "hadn": {}, "hadn't": {},
	"hasn": {}, "hasn't": {}, "haven": {}, "haven't": {}, "isn": {}, "isn't": {}, "ma": {},
	"mightn": {}, "mightn't": {}, "mustn": {}, "mustn't": {}, "needn": {}, "needn't": {}, "shan": {},
	"shan't": {}, "shouldn": {}, "shouldn't": {}, "wasn": {}, "wasn't": {}, "weren": {}, "weren't": {},
	"won": {}, "won't": {}, "wouldn": {}, "wouldn't": {},
	"http": {},
}

func RemoveStopWords(text string) string {
	words := strings.Fields(text)
	filteredWords := []string{}

	// Loop thru the words in the text and check for stop words
	for _, word := range words {
		if _, found := STOP_WORDS[word]; !found {
			filteredWords = append(filteredWords, word)
		}
	}

	return strings.Join(filteredWords, " ")
}

func NormalizeText(text string) string {
	// Convert text to lowercase
	text = strings.ToLower(text)

	// Remove punctuation and other non-word characters
	text = regexp.MustCompile(`[^\w\s]`).ReplaceAllString(text, "")

	// Remove extra whitespace
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return text
}

func StemText(text string) string {
	words := strings.Fields(text)
	stemmedWords := []string{}

	// Loop thru the words in the text and stem 'em
	for _, word := range words {
		stemmedWord, err := snowball.Stem(word, "english", true)
		if err != nil {
			stemmedWords = append(stemmedWords, word)
		} else {
			stemmedWords = append(stemmedWords, stemmedWord)
		}
	}

	return strings.Join(stemmedWords, " ")
}
