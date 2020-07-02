package cmd

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(learnCommand)
}

var learnCommand = &cobra.Command{
	Use:   "learn",
	Short: "Add a new word pair",
	Long:  "Adds a new word pair to your vocabulary",
	Run: func(cmd *cobra.Command, args []string) {
		vocabulary := MustVocabulary()

		pair, stats := vocabulary.getLeastConfidentWord()
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

var tags = learnCommand.Flags().StringArrayP("tag", "t", []string{}, "")

func (vocabulary Vocabulary) getLeastConfidentWord() (*WordPair, *WordStats) {
	bestScore := 999999
	index := 0
	for i, word := range vocabulary.Words {
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

	word := vocabulary.Words[index]
	return &word, vocabulary.Stats[word.Name]
}
