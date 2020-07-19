package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"

	voc "github.com/JaCoB1123/vocabulary/internal/vocabulary"
)

var statsTags *[]string

func init() {
	statsTags = statsCmd.Flags().StringArrayP("tag", "t", []string{}, "Include only words having all specified tags")
	rootCmd.AddCommand(statsCmd)
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show learning stats",
	Long:  "Show stats about current the learning progress",
	Run: func(cmd *cobra.Command, args []string) {
		vocabulary := voc.MustVocabulary(*wordsfilename, *statsfilename)

		levelStats := make([]int, 8)
		tagStats := map[string]int{}
		timeStats := map[time.Duration]int{}
		totalDue := 0
		totalAnswers := 0
		alwaysCorrect := 0
		wordsAnswered := 0
		for _, pair := range vocabulary.Words {
			if pair.IsFilteredBy(*statsTags) {
				continue
			}

			stats := vocabulary.GetStats(pair)
			totalAnswers = totalAnswers + stats.Answers
			if stats.CorrectAnswers > 0 && stats.FalseAnswers == 0 {
				alwaysCorrect++
			}

			if stats.Answers > 0 {
				wordsAnswered += 1
			}

			timeStats[stats.LastDuration()] = timeStats[stats.LastDuration()] + 1

			for _, tag := range pair.Tags {
				tagStats[tag] = tagStats[tag] + 1
			}
			levelStats[stats.AnswersSinceLastError] = levelStats[stats.AnswersSinceLastError] + 1

			if stats.IsDue() {
				totalDue += 1
			}
		}

		fmt.Printf("Words by number level:\n")
		for key, count := range levelStats {
			fmt.Printf("%4d: %8d words\n", key, count)
		}
		fmt.Println()

		timeKeys := make([]int, 0, len(timeStats))
		for k := range timeStats {
			timeKeys = append(timeKeys, int(k))
		}
		sort.Ints(timeKeys)

		fmt.Printf("Last answered:\n")
		for _, key := range timeKeys {
			duration := time.Duration(key)
			if duration == voc.MaxDuration {
				fmt.Printf("%22s: %8d words\n", "never", timeStats[duration])
			} else {
				fmt.Printf("up to %12s ago: %8d words\n", formatTime(duration), timeStats[duration])
			}
		}
		fmt.Println()

		tagKeys := make([]string, 0, len(tagStats))
		for k := range tagStats {
			tagKeys = append(tagKeys, k)
		}
		sort.Strings(tagKeys)

		fmt.Printf("Words by Tags:\n")
		for _, key := range tagKeys {
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

func formatTime(duration time.Duration) string {
	if duration < time.Hour {
		return fmt.Sprintf("%d minutes", int(duration.Minutes()))
	}

	if duration < time.Hour*24 {
		return fmt.Sprintf("%d hours", int(duration.Hours()))
	}

	if duration < time.Hour*24*7 {
		return fmt.Sprintf("%d days", int(duration.Hours()/24))
	}

	if duration < time.Hour*24*30 {
		return fmt.Sprintf("%d weeks", int(duration.Hours()/24/7))
	}

	return fmt.Sprintf("%d months", int(duration.Hours()/24/30))
}
