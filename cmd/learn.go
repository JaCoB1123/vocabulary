package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var learnTags *[]string

func init() {
	rootCmd.AddCommand(learnCommand)
	learnTags = learnCommand.Flags().StringArrayP("tag", "t", []string{}, "Include only words having all specified tags")
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
			stats.CorrectAnswer()
			fmt.Println("Correct!")
		} else {
			similarity := CompareTwoStrings(userInput, pair.Translation)
			fmt.Printf("Similarity: %f\n", similarity)

			acceptAnswer := similarity > 0.5

			fmt.Printf("Correct answer:\n%s\n\n", pair.Translation)

			acceptAnswer = promptTrueOrFalse(reader, "Accept Answer?", acceptAnswer)
			stats.CorrectAnswer()
			stats.FalseAnswer()
			fmt.Println("Wrong!")
		}

		vocabulary.Save()
	},
}

func (vocabulary Vocabulary) getLeastConfidentWord() (*WordPair, *WordStats, error) {
	bestScore := 999999
	index := -1
	for i, word := range vocabulary.Words {
		if !containsAll(*learnTags, word.Tags) {
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
		return nil, nil, fmt.Errorf("No word matching tags %v found", *learnTags)
	}

	word := vocabulary.Words[index]
	return &word, vocabulary.Stats[word.Name], nil
}
