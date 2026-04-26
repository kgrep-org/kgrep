package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/hbelmiro/kgrep/internal/resource"
)

// outputFormat controls how the output is displayed
var outputFormat string

// formatPatternForDisplay wraps the pattern in quotes if it contains special characters
// and is not already wrapped in quotes
func formatPatternForDisplay(pattern string) string {
	// If already wrapped in single or double quotes, return as-is
	if (strings.HasPrefix(pattern, "'") && strings.HasSuffix(pattern, "'")) ||
		(strings.HasPrefix(pattern, "\"") && strings.HasSuffix(pattern, "\"")) {
		return pattern
	}

	// Check if pattern contains special characters that would require quoting
	// These are characters that have special meaning in shells or could be confusing
	specialChars := []string{"[", "]", "*", "?", "{", "}", "|", "&", ";", "<", ">", "`", "$", "(", ")"}
	for _, char := range specialChars {
		if strings.Contains(pattern, char) {
			return fmt.Sprintf("\"%s\"", pattern)
		}
	}

	return pattern
}

// printResourceOccurrences prints occurrences based on the output format
func printResourceOccurrences(occurrences []resource.Occurrence, pattern string) {
	displayPattern := formatPatternForDisplay(pattern)

	if len(occurrences) == 0 {
		fmt.Printf("No occurrences of %s found.\n", displayPattern)
		return
	}

	if outputFormat == "name-only" {
		printResourceNamesOnly(occurrences)
		return
	}

	// Default format
	fmt.Printf("Found %d occurrence(s) of %s:\n\n", len(occurrences), displayPattern)

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
