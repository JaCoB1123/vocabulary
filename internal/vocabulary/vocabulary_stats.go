package vocabulary

import (
	"time"
)

type VocabularyStats struct {
	TimeStats     map[time.Duration]int
	LevelStats    []int
	TagStats      map[string]int
	TotalAnswers  int
	AlwaysCorrect int
	WordsAnswered int
	TotalDue      int
}

func (vocabulary Vocabulary) GetVocabularyStats(statsTags []string) *VocabularyStats {
	stats := &VocabularyStats{}
	stats.TimeStats = map[time.Duration]int{}
	stats.LevelStats = make([]int, 8)
	stats.TagStats = map[string]int{}
	for _, pair := range vocabulary.Words {
		if pair.IsFilteredBy(statsTags) {
			continue
		}

		wordStats := vocabulary.GetStats(pair)
		stats.TotalAnswers = stats.TotalAnswers + wordStats.Answers
		if wordStats.CorrectAnswers > 0 && wordStats.FalseAnswers == 0 {
			stats.AlwaysCorrect++
		}

		if wordStats.Answers > 0 {
			stats.WordsAnswered += 1
		}

		stats.TimeStats[wordStats.LastDuration()] = stats.TimeStats[wordStats.LastDuration()] + 1

		for _, tag := range pair.Tags {
			stats.TagStats[tag] = stats.TagStats[tag] + 1
		}
		stats.LevelStats[wordStats.AnswersSinceLastError] = stats.LevelStats[wordStats.AnswersSinceLastError] + 1

		if wordStats.IsDue() {
			stats.TotalDue += 1
		}
	}
	return stats
}
