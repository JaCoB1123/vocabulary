package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var importTags *[]string

var importFilename *string

var importSeparator *string

func init() {
	rootCmd.AddCommand(importCommand)
	importTags = importCommand.Flags().StringArrayP("tag", "t", []string{}, "Tags to assign to all imported word pairs")
	importFilename = importCommand.Flags().StringP("filename", "f", "", "File to import the word pairs from")
	importSeparator = importCommand.Flags().StringP("separator", "s", ",", "Separator to use instead of comma")
}

var importCommand = &cobra.Command{
	Use:   "import",
	Short: "Import word pairs",
	Long:  "Imports new word pairs from a comma separated values file",
	Run: func(cmd *cobra.Command, args []string) {
		vocabulary := MustVocabulary()

		file, err := os.Open(*importFilename)
		if err != nil {
			log.Fatalf("Could not open file: %s\n", err.Error())
		}

		reader := csv.NewReader(file)
		reader.Comma = ([]rune(*importSeparator))[0]
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			wordpair := WordPair{
				Name:        record[0],
				Translation: record[1],
				Tags:        *importTags,
			}

			vocabulary.Words = append(vocabulary.Words, wordpair)
			fmt.Println(wordpair)
		}

		vocabulary.Save()
	},
}
