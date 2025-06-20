package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/hbelmiro/kgrep/internal/resource"
)

func printResourceOccurrences(occurrences []resource.Occurrence, pattern string) {
	if len(occurrences) == 0 {
		fmt.Printf("No occurrences of '%s' found.\n", pattern)
		return
	}

	fmt.Printf("Found %d occurrence(s) of '%s':\n\n", len(occurrences), pattern)

	for _, occurrence := range occurrences {
		highlightedContent := strings.ReplaceAll(occurrence.Content, pattern, color.RedString(pattern))
		fmt.Printf("%s:%d: %s\n", occurrence.Resource, occurrence.Line, highlightedContent)
	}
}
