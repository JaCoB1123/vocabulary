package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCommand)
	tags = importCommand.Flags().StringArrayP("tag", "t", []string{}, "Tags to assign to all imported word pairs")
	filename = importCommand.Flags().StringP("filename", "f", "", "File to import the word pairs from")
}

var importCommand = &cobra.Command{
	Use:   "import",
	Short: "Import word pairs",
	Long:  "Imports new word pairs from a comma separated values file",
	Run: func(cmd *cobra.Command, args []string) {
		vocabulary := MustVocabulary()

		vocabulary.Save()
	},
}
