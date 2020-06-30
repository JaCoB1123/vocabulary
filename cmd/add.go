package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCommand)
}

var addCommand = &cobra.Command{
	Use:   "add [name] [translation] [optional tags...]",
	Short: "Add a new word pair",
	Long:  "Adds a new word pair to your vocabulary",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vocabulary := MustVocabulary()

		vocabulary.Words = append(vocabulary.Words, WordPair{
			Name:        args[0],
			Translation: args[1],
			Tags:        args[2:],
		})

		vocabulary.Save()
	},
}
