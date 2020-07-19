package vocabulary

import (
	"math"
	"time"
)

type WordStats struct {
	Answers               int       `json:",omitempty"`
	CorrectAnswers        int       `json:",omitempty"`
	FalseAnswers          int       `json:",omitempty"`
	LastCorrect           time.Time `json:",omitempty"`
	LastFalse             time.Time `json:",omitempty"`
	AnswersSinceLastError int       `json:",omitempty"`
}

func (stats *WordStats) CorrectAnswer() {
	stats.AnswersSinceLastError++
	stats.CorrectAnswers++
	stats.LastCorrect = time.Now()
}

func (stats *WordStats) FalseAnswer() {
	stats.AnswersSinceLastError = 0
	stats.FalseAnswers++
	stats.LastFalse = time.Now()
}

func (stats *WordStats) LastAnswered() time.Time {
	if stats.LastCorrect.After(stats.LastFalse) {
		return stats.LastCorrect
	}

	return stats.LastFalse
}

func (stats WordStats) RecommendedDuration() time.Duration {
	return getRecommendedDuration(stats.AnswersSinceLastError)
}

var MaxDuration = time.Duration(math.MaxInt64)

func (stats WordStats) LastDuration() time.Duration {
	if stats.LastAnswered().IsZero() {
		return MaxDuration
	}
	return getRecommendedDuration(getRecommendedScore(time.Now().Sub(stats.LastAnswered())))
}

// check if word should be practiced
func (stats WordStats) IsDue() bool {
	if stats.LastAnswered().IsZero() {
		return true
	}

	requiredAge := stats.RecommendedDuration()
	dueOn := stats.LastAnswered().Add(requiredAge)
	return time.Now().After(dueOn)
}

func (stats WordStats) GetScore() int64 {
	if !stats.IsDue() {
		return math.MinInt64
	}

	score := int64(stats.AnswersSinceLastError + 1)

	if !stats.LastAnswered().IsZero() {
		score = score * int64(getRecommendedScore(time.Now().Sub(stats.LastAnswered())))
	}
	return score
}

var recommendedDurations = []time.Duration{
	time.Minute * 10,
	time.Hour * 1,
	time.Hour * 4,
	time.Hour * 24,
	time.Hour * 24 * 7,
	time.Hour * 24 * 30,
	time.Hour * 24 * 30 * 6,
}

func getRecommendedDuration(sucessfullTries int) time.Duration {
	if sucessfullTries < 0 {
		return recommendedDurations[0]
	}

	if sucessfullTries > len(recommendedDurations) {
		return recommendedDurations[len(recommendedDurations)-1]
	}

	return recommendedDurations[sucessfullTries]
}

func getRecommendedScore(duration time.Duration) int {
	for index, recommendedDuration := range recommendedDurations {
		if duration < recommendedDuration {
			return index
		}
	}

	return len(recommendedDurations) - 1
}
