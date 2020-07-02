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
