package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(learnCommand)
	tags = learnCommand.Flags().StringArrayP("tag", "t", []string{}, "Include only words having all specified tags")
}

var learnCommand = &cobra.Command{
	Use:   "learn",
	Short: "Add a new word pair",
	Long:  "Adds a new word pair to your vocabulary",
	Run: func(cmd *cobra.Command, args []string) {
		vocabulary := MustVocabulary()

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
			stats.AnswersSinceLastError++
			stats.CorrectAnswers++
			stats.LastCorrect = time.Now()

			fmt.Println("Correct!")
		} else {
			stats.AnswersSinceLastError = 0
			stats.FalseAnswers++
			stats.LastFalse = time.Now()
			fmt.Println("Wrong!")
			fmt.Printf("Your answer: '%v'\n", []byte(userInput))
			fmt.Printf("Correct answer: '%s'\n", pair.Translation)
		}

		vocabulary.Save()
	},
}

var tags *[]string

func (vocabulary Vocabulary) getLeastConfidentWord() (*WordPair, *WordStats, error) {
	bestScore := 999999
	index := -1
	for i, word := range vocabulary.Words {
		if !containsAll(*tags, word.Tags) {
			continue
		}

		var stats *WordStats
		if mapstats, ok := vocabulary.Stats[word.Name]; ok {
			stats = mapstats
		} else {
			stats = &WordStats{}
			vocabulary.Stats[word.Name] = stats
		}

		score := stats.AnswersSinceLastError
		if score < bestScore {
			index = i
			bestScore = score
		}
	}

	if index == -1 {
		return nil, nil, fmt.Errorf("No word matching tags %v found", *tags)
	}

	word := vocabulary.Words[index]
	return &word, vocabulary.Stats[word.Name], nil
}
