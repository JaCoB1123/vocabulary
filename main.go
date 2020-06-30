package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"encoding/json"
)

type Vocabulary struct {
	Words []WordPair
	Stats map[string]WordStats
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

func main() {
	vocabulary, err := NewVocabulary()
	if err != nil {
		log.Fatalf("Error loading vocabulary: %s", err.Error())
	}

	fmt.Printf("%v\n", vocabulary)
}

func NewVocabulary() (*Vocabulary, error) {
	var vocabulary Vocabulary

	wordsfilename := "words.json"
	statsfilename := "stats.json"

	var words []WordPair
	err := DeserializeFile(wordsfilename, &words)
	if err != nil {
		return nil, fmt.Errorf("Could not read words: %v", err)
	}
	vocabulary.Words = words

	var stats map[string]WordStats
	err = DeserializeFile(statsfilename, &stats)
	if err != nil {
		return nil, fmt.Errorf("Could not read stats %v", err)
	}
	vocabulary.Stats = stats

	return &vocabulary, nil
}

func DeserializeFile(filename string, target interface{}) error {
	wordsfile, err := os.Open(filename)
	if os.IsNotExist(err) {
		wordsfile, err = os.Create(filename)
		if err != nil {
			return fmt.Errorf("Could not create file %s: %v", filename, err)
		}

		err = json.NewEncoder(wordsfile).Encode(target)
		if err != nil {
			return fmt.Errorf("Could not write JSON to file %s: %v", filename, err)
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("Could not open file %s: %v", filename, err)
	}

	var words []WordPair
	err = json.NewDecoder(wordsfile).Decode(&words)
	if err != nil {
		return fmt.Errorf("Could not deserialize JSON %s: %v", filename, err)
	}

	return nil
}
