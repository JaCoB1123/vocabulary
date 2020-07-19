package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/agnivade/levenshtein"
	"github.com/spf13/cobra"

	voc "github.com/JaCoB1123/vocabulary/internal/vocabulary"
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
		vocabulary := voc.MustVocabulary(*wordsfilename, *statsfilename)

		for i := 0; i < *count; i++ {
			pair, stats, err := vocabulary.GetLeastConfidentWord(*learnTags)
			if err != nil {
				log.Printf("No words found: %s\n", err.Error())
				break
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

		vocabulary.Save(*wordsfilename, *statsfilename)
	},
}
