package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vocabulary",
	Short: "Vocabulary is a simple vocabulary management and practicing tool",
	Long:  `A simple vocabulary management and practicing tool for the command line`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var wordsfilename *string
var statsfilename *string

func Execute() {
	statsfilename = rootCmd.PersistentFlags().String("stats-file", "stats.json", "Loads stats from the specified file")
	wordsfilename = rootCmd.PersistentFlags().String("words-file", "words.json", "Loads words from the specified file")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
