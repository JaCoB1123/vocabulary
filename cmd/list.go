package cmd

import (
	"fmt"

	voc "github.com/JaCoB1123/vocabulary/internal/vocabulary"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all word pairs",
	Long:  "Used to print all word pairs in different formats",
	Run: func(cmd *cobra.Command, args []string) {
		vocabulary := voc.MustVocabulary(*wordsfilename, *statsfilename)

		for _, pair := range vocabulary.Words {
			stats := vocabulary.GetStats(pair)
			score := stats.GetScore()

			fmt.Printf("%40s: %40s (%10d)\n", pair.Name, pair.Translation, score)
		}
	},
}
