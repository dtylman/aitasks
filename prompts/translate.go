package prompts

import "fmt"

// GetTranslateStyles returns a list of available styles for the translate task by reading the embedded prompt templates.
func GetTranslateStyles() ([]string, error) {
	entries, err := embeddedFS.ReadDir("embedded/translate")
	if err != nil {
		return nil, fmt.Errorf("read translate styles: %w", err)
	}
	var styles []string
	for _, entry := range entries {
		if entry.IsDir() {
			styles = append(styles, entry.Name())
		}
	}
	return styles, nil
}
