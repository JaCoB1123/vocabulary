package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statsCmd)
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show learning stats",
	Long:  "Show stats about current the learning progress",
	Run: func(cmd *cobra.Command, args []string) {
		vocabulary := MustVocabulary()

		levelStats := make([]int, 8)
		tagStats := map[string]int{}
		totalDue := 0
		totalAnswers := 0
		alwaysCorrect := 0
		wordsAnswered := 0
		for _, pair := range vocabulary.Words {
			stats := vocabulary.getStats(pair)
			totalAnswers = totalAnswers + stats.Answers
			if stats.CorrectAnswers > 0 && stats.FalseAnswers == 0 {
				alwaysCorrect++
			}

			if stats.Answers > 0 {
				wordsAnswered += 1
			}

			for _, tag := range pair.Tags {
				tagStats[tag] = tagStats[tag] + 1
			}
			levelStats[stats.AnswersSinceLastError] = levelStats[stats.AnswersSinceLastError] + 1

			if stats.isDue() {
				totalDue += 1
			}
		}

		fmt.Printf("Words by number level:\n")
		for key, count := range levelStats {
			fmt.Printf("%4d: %8d words\n", key, count)
		}
		fmt.Println()

		keys := make([]string, 0, len(tagStats))
		for k := range tagStats {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		fmt.Printf("Words by Tags:\n")
		for _, key := range keys {
			fmt.Printf("%20s: %8d words\n", key, tagStats[key])
		}
		fmt.Println()

		fmt.Printf("Words answered at least once: %7d\n", wordsAnswered)
		fmt.Printf("Words with 100%% success:      %7d\n", alwaysCorrect)
		fmt.Printf("words currently due:          %7d\n", totalDue)

		fmt.Printf("Total words:                  %7d\n", len(vocabulary.Words))
		fmt.Printf("Total answers:                %7d\n", totalAnswers)
	},
}
