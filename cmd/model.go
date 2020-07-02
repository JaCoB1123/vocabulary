package cmd

import (
	"time"
)

type Vocabulary struct {
	Words []WordPair
	Stats map[string]*WordStats
}

type WordStats struct {
	Answers               int
	CorrectAnswers        int
	FalseAnswers          int
	LastCorrect           time.Time
	LastFalse             time.Time
	AnswersSinceLastError int
}

type WordPair struct {
	Name        string
	Translation string
	Tags        []string
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
