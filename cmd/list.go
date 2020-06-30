package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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
		vocabulary := MustVocabulary()

		for _, pair := range vocabulary.Words {
			fmt.Printf("%s: %s\n", pair.Name, pair.Translation)
		}
	},
}

var wordsfilename = "words.json"
var statsfilename = "stats.json"

func (vocabulary Vocabulary) Save() {
	SerializeFile(wordsfilename, vocabulary.Words)
	SerializeFile(statsfilename, vocabulary.Stats)
}

func MustVocabulary() *Vocabulary {
	vocabulary, err := NewVocabulary()
	if err != nil {
		log.Fatalf("Error loading vocabulary: %s", err.Error())
	}

	return vocabulary
}

func NewVocabulary() (*Vocabulary, error) {
	var vocabulary Vocabulary

	var words []WordPair
	err := DeserializeFile(wordsfilename, &words)
	if err != nil {
		return nil, fmt.Errorf("Could not read words: %v", err)
	}
	vocabulary.Words = words

	var stats map[string]*WordStats
	err = DeserializeFile(statsfilename, &stats)
	if err != nil {
		return nil, fmt.Errorf("Could not read stats %v", err)
	}
	vocabulary.Stats = stats

	return &vocabulary, nil
}

func SerializeFile(filename string, target interface{}) error {
	wordsfile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Could not create file %s: %v", filename, err)
	}

	encoder := json.NewEncoder(wordsfile)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(target)
	if err != nil {
		return fmt.Errorf("Could not write JSON to file %s: %v", filename, err)
	}

	return nil
}

func DeserializeFile(filename string, target interface{}) error {
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return SerializeFile(filename, target)
	}

	if err != nil {
		return fmt.Errorf("Could not open file %s: %v", filename, err)
	}

	err = json.NewDecoder(file).Decode(&target)
	if err != nil {
		return fmt.Errorf("Could not deserialize JSON %s: %v", filename, err)
	}

	return nil
}
