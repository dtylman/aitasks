package prompts

import (
	"reflect"
	"testing"
)

func TestGetTranslateStyles_ReturnsExpectedStyles(t *testing.T) {
	styles, err := GetTranslateStyles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"academic", "default", "literary", "strict"}
	if !reflect.DeepEqual(styles, expected) {
		t.Fatalf("expected styles %v, got %v", expected, styles)
	}
}

func TestGetTranslateStyles_MatchesEmbeddedTranslateDirectories(t *testing.T) {
	entries, err := embeddedFS.ReadDir("embedded/translate")
	if err != nil {
		t.Fatalf("failed to read embedded translate directory: %v", err)
	}

	var expected []string
	for _, entry := range entries {
		if entry.IsDir() {
			expected = append(expected, entry.Name())
		}
	}

	styles, err := GetTranslateStyles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(styles, expected) {
		t.Fatalf("expected styles %v, got %v", expected, styles)
	}
}
