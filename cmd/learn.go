package cmd

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/agnivade/levenshtein"
	"github.com/spf13/cobra"
)

var learnTags *[]string

var count *int

func init() {
	rootCmd.AddCommand(learnCommand)
	learnTags = learnCommand.Flags().StringArrayP("tag", "t", []string{}, "Include only words having all specified tags")
	count = learnCommand.Flags().IntP("count", "n", 1, "Ask this number of words before quitting")
}

var learnCommand = &cobra.Command{
	Use:   "learn",
	Short: "Add a new word pair",
	Long:  "Adds a new word pair to your vocabulary",
	Run: func(cmd *cobra.Command, args []string) {
		vocabulary := MustVocabulary()

		for i := 0; i < *count; i++ {
			pair, stats, err := vocabulary.getLeastConfidentWord()
			if err != nil {
				log.Fatalf("Error finding word: %s\n", err.Error())
			}
			fmt.Println("Word: ", pair.Name)

			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Answer: ")
			userInput, _ := reader.ReadString('\n')
			userInput = trimSuffixLineEnding(userInput)
			fmt.Println(userInput)

			stats.Answers++
			if userInput == pair.Translation {
				stats.CorrectAnswer()
				fmt.Println("Correct!")
			} else {
				distance := levenshtein.ComputeDistance(userInput, pair.Translation)
				length := len(userInput)
				if len(pair.Translation) > length {
					length = len(pair.Translation)
				}
				similarity := 1 - float32(distance)/float32(length)
				fmt.Printf("Similarity: %f\n", similarity)

				acceptAnswer := similarity > 0.5

				fmt.Printf("Correct answer:\n%s\n\n", pair.Translation)

				acceptAnswer = promptTrueOrFalse(reader, "Accept Answer?", acceptAnswer)
				if acceptAnswer {
					stats.CorrectAnswer()
				} else {
					stats.FalseAnswer()
				}
			}
		}

		vocabulary.Save()
	},
}

func (vocabulary Vocabulary) getStats(word WordPair) *WordStats {
	var stats *WordStats
	if mapstats, ok := vocabulary.Stats[word.Name]; ok {
		return mapstats
	}
	stats = &WordStats{}
	vocabulary.Stats[word.Name] = stats
	return stats

}

// check if word should be practiced again
func (stats WordStats) isDue() bool {
	requiredAge := getRecommendedDuration(stats.AnswersSinceLastError)
	return stats.LastCorrect.Before(time.Now().Add(requiredAge))
}

func (stats WordStats) getScore() int64 {
	if !stats.isDue() {
		return math.MinInt64
	}

	score := int64(stats.AnswersSinceLastError + 1)

	if !stats.LastAnswered().IsZero() {
		score = score * stats.LastAnswered().Unix()
	}
	return score
}

func (vocabulary Vocabulary) getLeastConfidentWord() (*WordPair, *WordStats, error) {
	bestScore := int64(math.MaxInt64)
	index := -1
	for i, word := range vocabulary.Words {
		if !containsAll(*learnTags, word.Tags) {
			continue
		}

		stats := vocabulary.getStats(word)

		score := stats.getScore()
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
		return nil, nil, fmt.Errorf("No word matching tags %v found", *learnTags)
	}

	word := vocabulary.Words[index]
	return &word, vocabulary.Stats[word.Name], nil
}

func getRecommendedDuration(sucessfullTries int) time.Duration {
	if sucessfullTries <= 0 {
		return time.Duration(0)
	}

	switch sucessfullTries {
	case 1:
		return time.Duration(time.Minute * 10)
	case 2:
		return time.Duration(time.Hour)
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
