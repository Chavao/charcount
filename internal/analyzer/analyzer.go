package analyzer

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	wordPattern      = regexp.MustCompile(`[\p{L}][\p{L}\p{N}'-]*`)
	sentencePattern  = regexp.MustCompile(`[^.!?]+[.!?]+|[^.!?]+$`)
	paragraphPattern = regexp.MustCompile(`\n\s*\n`)
)

type DensityRow struct {
	Word    string `json:"word"`
	Count   int    `json:"count"`
	Density string `json:"density"`
}

type Analysis struct {
	Characters  int          `json:"characters"`
	Words       int          `json:"words"`
	Sentences   int          `json:"sentences"`
	Paragraphs  int          `json:"paragraphs"`
	Spaces      int          `json:"spaces"`
	DensityRows []DensityRow `json:"densityRows"`
}

func Analyze(text string) Analysis {
	words := extractWords(text)
	wordCount := len(words)
	rows := buildDensityRows(words, wordCount)

	return Analysis{
		Characters:  len(text),
		Words:       wordCount,
		Sentences:   countSentences(text),
		Paragraphs:  countParagraphs(text),
		Spaces:      strings.Count(text, " "),
		DensityRows: rows,
	}
}

func extractWords(text string) []string {
	matches := wordPattern.FindAllString(text, -1)
	if len(matches) == 0 {
		return nil
	}

	words := make([]string, 0, len(matches))
	for _, match := range matches {
		words = append(words, strings.ToLower(match))
	}

	return words
}

func countSentences(text string) int {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return 0
	}

	matches := sentencePattern.FindAllString(trimmed, -1)
	if len(matches) == 0 {
		return 0
	}

	sentences := 0
	for _, sentence := range matches {
		if strings.TrimSpace(sentence) != "" {
			sentences++
		}
	}

	return sentences
}

func countParagraphs(text string) int {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	normalized = strings.TrimSpace(normalized)
	if normalized == "" {
		return 0
	}

	paragraphs := paragraphPattern.Split(normalized, -1)
	count := 0
	for _, paragraph := range paragraphs {
		if strings.TrimSpace(paragraph) != "" {
			count++
		}
	}

	return count
}

func buildDensityRows(words []string, wordCount int) []DensityRow {
	if wordCount == 0 {
		return []DensityRow{}
	}

	counts := make(map[string]int, len(words))
	for _, word := range words {
		counts[word]++
	}

	sortedWords := make([]string, 0, len(counts))
	for word := range counts {
		sortedWords = append(sortedWords, word)
	}

	sort.Slice(sortedWords, func(i int, j int) bool {
		left := sortedWords[i]
		right := sortedWords[j]

		leftCount := counts[left]
		rightCount := counts[right]
		if leftCount != rightCount {
			return leftCount > rightCount
		}

		return left < right
	})

	rows := make([]DensityRow, 0, len(sortedWords))
	for _, word := range sortedWords {
		count := counts[word]
		density := fmt.Sprintf("%.2f%%", float64(count)/float64(wordCount)*100)
		rows = append(rows, DensityRow{
			Word:    word,
			Count:   count,
			Density: density,
		})
	}

	return rows
}
