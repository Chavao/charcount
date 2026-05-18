package analyzer_test

import (
	"testing"

	"github.com/Chavao/charcount/internal/analyzer"
)

func TestAnalyze(t *testing.T) {
	t.Parallel()

	text := "Hello world! Hello, Go.\n\nNext paragraph"

	got := analyzer.Analyze(text)

	if got.Characters != len(text) {
		t.Fatalf("Characters = %d, want %d", got.Characters, len(text))
	}

	if got.Words != 6 {
		t.Fatalf("Words = %d, want 6", got.Words)
	}

	if got.Sentences != 3 {
		t.Fatalf("Sentences = %d, want 3", got.Sentences)
	}

	if got.Paragraphs != 2 {
		t.Fatalf("Paragraphs = %d, want 2", got.Paragraphs)
	}

	if got.Spaces != 4 {
		t.Fatalf("Spaces = %d, want 4", got.Spaces)
	}

	wantRows := []analyzer.DensityRow{
		{Word: "hello", Count: 2, Density: "33.33%"},
		{Word: "go", Count: 1, Density: "16.67%"},
		{Word: "next", Count: 1, Density: "16.67%"},
		{Word: "paragraph", Count: 1, Density: "16.67%"},
		{Word: "world", Count: 1, Density: "16.67%"},
	}

	if len(got.DensityRows) != len(wantRows) {
		t.Fatalf("DensityRows length = %d, want %d", len(got.DensityRows), len(wantRows))
	}

	for i := range wantRows {
		if got.DensityRows[i] != wantRows[i] {
			t.Fatalf("DensityRows[%d] = %+v, want %+v", i, got.DensityRows[i], wantRows[i])
		}
	}
}

func TestAnalyzeEmptyText(t *testing.T) {
	t.Parallel()

	got := analyzer.Analyze("")

	if got.Characters != 0 {
		t.Fatalf("Characters = %d, want 0", got.Characters)
	}

	if got.Words != 0 {
		t.Fatalf("Words = %d, want 0", got.Words)
	}

	if got.Sentences != 0 {
		t.Fatalf("Sentences = %d, want 0", got.Sentences)
	}

	if got.Paragraphs != 0 {
		t.Fatalf("Paragraphs = %d, want 0", got.Paragraphs)
	}

	if got.Spaces != 0 {
		t.Fatalf("Spaces = %d, want 0", got.Spaces)
	}

	if len(got.DensityRows) != 0 {
		t.Fatalf("DensityRows length = %d, want 0", len(got.DensityRows))
	}
}
