package cmd

import (
	"time"
)

type Vocabulary struct {
	Words []WordPair
	Stats map[string]*WordStats
}

type WordStats struct {
	Answers               int       `json:",omitempty"`
	CorrectAnswers        int       `json:",omitempty"`
	FalseAnswers          int       `json:",omitempty"`
	LastCorrect           time.Time `json:",omitempty"`
	LastFalse             time.Time `json:",omitempty"`
	AnswersSinceLastError int       `json:",omitempty"`
}

type WordPair struct {
	Name        string
	Translation string
	Attributes  map[string]string `json:",omitempty"`
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

func (stats *WordStats) LastAnswered() time.Time {
	if stats.LastCorrect.After(stats.LastFalse) {
		return stats.LastCorrect
	}

	return stats.LastFalse
}
