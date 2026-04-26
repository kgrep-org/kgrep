package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/hbelmiro/kgrep/internal/resource"
)

// outputFormat controls how the output is displayed
var outputFormat string

// printResourceOccurrences prints occurrences based on the output format
func printResourceOccurrences(occurrences []resource.Occurrence, pattern string) {
	if len(occurrences) == 0 {
		fmt.Printf("No occurrences of '%s' found.\n", pattern)
		return
	}

	if outputFormat == "name-only" {
		printResourceNamesOnly(occurrences)
		return
	}

	// Default format
	fmt.Printf("Found %d occurrence(s) of '%s':\n\n", len(occurrences), pattern)

	for _, occurrence := range occurrences {
		boldRed := color.New(color.FgRed).Add(color.Bold)

		highlightedContent := strings.ReplaceAll(occurrence.Content, pattern, boldRed.Sprint(pattern))

		var prefix string
		if occurrence.Namespace != "" {
			prefix = color.BlueString("%s/%s[%d]:", occurrence.Namespace, occurrence.Resource, occurrence.Line)
		} else {
			prefix = color.BlueString("%s[%d]:", occurrence.Resource, occurrence.Line)
		}

		fmt.Printf("%s %s\n", prefix, highlightedContent)
	}
}

// printResourceNamesOnly prints only the unique resource names
func printResourceNamesOnly(occurrences []resource.Occurrence) {
	// Use a map to deduplicate resource names
	resourceNames := make(map[string]bool)
	for _, occurrence := range occurrences {
		var resourceKey string
		if occurrence.Namespace != "" {
			resourceKey = fmt.Sprintf("%s/%s", occurrence.Namespace, occurrence.Resource)
		} else {
			resourceKey = occurrence.Resource
		}
		resourceNames[resourceKey] = true
	}

	// Print unique resource names
	for name := range resourceNames {
		fmt.Println(name)
	}
}
