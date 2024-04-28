package hw03frequencyanalysis

import (
	"cmp"
	"slices"
	"strings"
)

func Top10(str string) []string {
	words := strings.Fields(str)

	if len(words) == 0 {
		return []string{}
	}

	wordsFreqMap := make(map[string]int)

	for _, word := range words {
		wordsFreqMap[word]++
	}

	words = make([]string, 0, len(wordsFreqMap))

	for value := range wordsFreqMap {
		words = append(words, value)
	}

	slices.SortFunc(words, func(a string, b string) int {
		if wordsFreqMap[b] == wordsFreqMap[a] {
			return cmp.Compare(a, b)
		}
		return cmp.Compare(wordsFreqMap[b], wordsFreqMap[a])
	})

	if len(words) > 10 {
		return words[:10]
	}

	return words
}
