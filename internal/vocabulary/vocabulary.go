package vocabulary

import (
	"fmt"
	"log"
	"math"

	"github.com/JaCoB1123/vocabulary/internal/serialization"
)

type Vocabulary struct {
	Words []WordPair
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

func (vocabulary Vocabulary) GetStats(word WordPair) *WordStats {
	var stats *WordStats
	if mapstats, ok := vocabulary.Stats[word.Name]; ok {
		return mapstats
	}
	stats = &WordStats{}
	vocabulary.Stats[word.Name] = stats
	return stats

}

func (vocabulary Vocabulary) GetLeastConfidentWord(learnTags []string) (*WordPair, *WordStats, error) {
	bestScore := int64(math.MaxInt64)
	index := -1
	for i, word := range vocabulary.Words {
		if !containsAll(learnTags, word.Tags) {
			continue
		}

		stats := vocabulary.GetStats(word)

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
		return nil, nil, fmt.Errorf("No word matching tags %v found", learnTags)
	}

	word := vocabulary.Words[index]
	return &word, vocabulary.Stats[word.Name], nil
}
