package vocabulary

import (
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/JaCoB1123/vocabulary/internal/serialization"
)

type Vocabulary struct {
	Words WordsList
	Stats map[string]*WordStats
}

func MustVocabulary(wordsfilename, statsfilename string) *Vocabulary {
	vocabulary, err := NewVocabulary(wordsfilename, statsfilename)
	if err != nil {
		log.Fatalf("Error loading vocabulary: %s", err.Error())
	}

	return vocabulary
}

func NewVocabulary(wordsfilename, statsfilename string) (*Vocabulary, error) {
	var vocabulary Vocabulary

	var words []WordPair
	err := serialization.DeserializeFile(wordsfilename, &words)
	if err != nil {
		return nil, fmt.Errorf("Could not read words: %v", err)
	}
	vocabulary.Words = words

	var stats map[string]*WordStats
	err = serialization.DeserializeFile(statsfilename, &stats)
	if err != nil {
		return nil, fmt.Errorf("Could not read stats %v", err)
	}
	vocabulary.Stats = stats

	if vocabulary.Stats == nil {
		vocabulary.Stats = make(map[string]*WordStats)
	}

	return &vocabulary, nil
}

func (vocabulary Vocabulary) Save(wordsfilename, statsfilename string) {
	serialization.SerializeFile(wordsfilename, vocabulary.Words)
	serialization.SerializeFile(statsfilename, vocabulary.Stats)
}

func (vocabulary Vocabulary) GetStatsByWord(word WordPair) *WordStats {
	return vocabulary.GetStats(word.Name)
}

func (vocabulary Vocabulary) GetStats(word string) *WordStats {
	var stats *WordStats
	if mapstats, ok := vocabulary.Stats[word]; ok {
		return mapstats
	}
	stats = &WordStats{}
	vocabulary.Stats[word] = stats
	return stats
}

func (vocabulary Vocabulary) GetSortedWords(tags []string) []WordPair {
	words := vocabulary.Words
	words = words.FilterByTags(tags)
	words = words.FilterRecent(vocabulary)

	sort.Slice(words, func(i, j int) bool {
		wordi := words[i]
		statsi := vocabulary.GetStatsByWord(wordi)
		scorei := statsi.GetScore()

		wordj := words[j]
		statsj := vocabulary.GetStatsByWord(wordj)
		scorej := statsj.GetScore()

		return scorei < scorej
	})

	return words
}

func (vocabulary Vocabulary) GetLeastConfidentWord(tags []string) (*WordPair, *WordStats, error) {
	bestScore := int64(math.MaxInt64)
	index := -1
	for i, word := range vocabulary.Words {
		if word.IsFilteredBy(tags) {
			continue
		}

		stats := vocabulary.GetStatsByWord(word)

		score := stats.GetScore()
		if score == math.MinInt64 {
			// ignore word as it has been answered recently
			continue
		}

		if score < bestScore {
			index = i
			bestScore = score
		}
	}

	if index == -1 {
		return nil, nil, fmt.Errorf("No word matching tags %v found", tags)
	}

	word := vocabulary.Words[index]
	return &word, vocabulary.Stats[word.Name], nil
}
