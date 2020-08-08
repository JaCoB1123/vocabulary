package vocabulary

import (
	"math"
)

type WordsList []WordPair

func (words WordsList) FilterByTags(tags []string) WordsList {
	if len(tags) == 0 {
		return words
	}

	filteredWords := []WordPair{}

	for _, word := range words {
		if word.IsFilteredBy(tags) {
			continue
		}

		filteredWords = append(filteredWords, word)
	}

	return filteredWords
}

func (words WordsList) FilterRecent(vocabulary Vocabulary) WordsList {
	filteredWords := []WordPair{}

	for _, word := range words {
		stats := vocabulary.GetStatsByWord(word)
		score := stats.GetScore()
		if score == math.MinInt64 {
			continue
		}

		filteredWords = append(filteredWords, word)
	}

	return filteredWords
}
