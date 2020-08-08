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

		stats := vocabulary.GetVocabularyStats(*statsTags)

		fmt.Printf("Words by number level:\n")
		for key, count := range stats.LevelStats {
			fmt.Printf("%4d: %8d words\n", key, count)
		}
		fmt.Println()

		timeKeys := make([]int, 0, len(stats.TimeStats))
		for k := range stats.TimeStats {
			timeKeys = append(timeKeys, int(k))
		}
		sort.Ints(timeKeys)

		fmt.Printf("Last answered:\n")
		for _, key := range timeKeys {
			duration := time.Duration(key)
			if duration == voc.MaxDuration {
				fmt.Printf("%22s: %8d words\n", "never", stats.TimeStats[duration])
			} else {
				fmt.Printf("up to %12s ago: %8d words\n", formatTime(duration), stats.TimeStats[duration])
			}
		}
		fmt.Println()

		tagKeys := make([]string, 0, len(stats.TagStats))
		for k := range stats.TagStats {
			tagKeys = append(tagKeys, k)
		}
		sort.Strings(tagKeys)

		fmt.Printf("Words by Tags:\n")
		for _, key := range tagKeys {
			fmt.Printf("%20s: %8d words\n", key, stats.TagStats[key])
		}
		fmt.Println()

		fmt.Printf("Words answered at least once: %7d\n", stats.WordsAnswered)
		fmt.Printf("Words with 100%% success:      %7d\n", stats.AlwaysCorrect)
		fmt.Printf("words currently due:          %7d\n", stats.TotalDue)

		fmt.Printf("Total words:                  %7d\n", stats.TotalWords)
		fmt.Printf("Total answers:                %7d\n", stats.TotalAnswers)
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
