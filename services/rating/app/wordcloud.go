package app

import (
	"strings"
	"unicode"
)

func BuildWordCloud(text string) map[string]int {
	result := make(map[string]int)

	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	stopWords := map[string]bool{
		"и": true, "в": true, "на": true, "с": true,
		"the": true, "a": true, "to": true,
	}

	for _, w := range words {
		if len(w) < 3 || stopWords[w] {
			continue
		}
		result[w]++
	}

	return result
}
