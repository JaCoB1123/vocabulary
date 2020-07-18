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

// check if word should be practiced
func (stats WordStats) IsDue() bool {
	if stats.LastCorrect.IsZero() {
		return true
	}

	requiredAge := getRecommendedDuration(stats.AnswersSinceLastError)
	dueOn := stats.LastCorrect.Add(requiredAge)
	return time.Now().After(dueOn)
}

func (stats WordStats) GetScore() int64 {
	if !stats.IsDue() {
		return math.MinInt64
	}

	score := int64(stats.AnswersSinceLastError + 1)

	if !stats.LastAnswered().IsZero() {
		score = score * stats.LastAnswered().Unix()
	}
	return score
}

func getRecommendedDuration(sucessfullTries int) time.Duration {
	if sucessfullTries <= 0 {
		return time.Duration(0)
	}

	switch sucessfullTries {
	case 1:
		return time.Duration(time.Minute * 30)
	case 2:
		return time.Duration(time.Hour * 3)
	case 3:
		return time.Duration(time.Hour * 24)
	case 4:
		return time.Duration(time.Hour * 24 * 7)
	case 5:
		return time.Duration(time.Hour * 24 * 30)
	default:
		return time.Duration(time.Hour * 24 * 30 * 6)
	}
}
