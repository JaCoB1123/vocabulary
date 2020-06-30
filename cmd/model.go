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
