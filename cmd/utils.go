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
		boldRed := color.New(color.FgRed).Add(color.Bold)

		highlightedContent := strings.ReplaceAll(occurrence.Content, pattern, boldRed.Sprint(pattern))
		prefix := color.BlueString("%s[%d]:", occurrence.Resource, occurrence.Line)
		fmt.Printf("%s %s\n", prefix, highlightedContent)
	}
}
