package cmd

import (
	"strings"
)

func trimSuffixLineEnding(text string) string {
	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")
	text = strings.TrimSuffix(text, "\n")
	return text
}
