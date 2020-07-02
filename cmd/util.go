package cmd

import (
	"bufio"
	"fmt"
	"strings"
)

func trimSuffixLineEnding(text string) string {
	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")
	text = strings.TrimSuffix(text, "\n")
	return text
}

func containsAll(needles []string, list []string) bool {
	for _, needle := range needles {
		if !contains(needle, list) {
			return false
		}
	}
	return true
}

func contains(needle string, list []string) bool {
	for _, b := range list {
		if b == needle {
			return true
		}
	}

	return false
}

func promptTrueOrFalse(reader *bufio.Reader, question string, fallback bool) bool {
	fmt.Printf("%s (true/false, default: %t) ", question, fallback)
	userInput, _ := reader.ReadString('\n')
	userInput = trimSuffixLineEnding(userInput)

	switch userInput {
	case "true":
		fallthrough
	case "yes":
		fallthrough
	case "y":
		return true
	case "false":
		fallthrough
	case "no":
		fallthrough
	case "n":
		return false
	case "":
		return fallback
	default:
		fmt.Println("Invalid input, please type true/yes/y or false/no/n")
		return promptTrueOrFalse(reader, question, fallback)
	}
}
