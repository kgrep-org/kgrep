package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/hbelmiro/kgrep/internal/log"
	"github.com/spf13/cobra"
)

var (
	logsNamespace string
	logsResource  string
	logsPattern   string
	logsSortBy    string
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Search logs in Kubernetes",
	Long:  `Search logs from a group of pods or entire namespaces, filtering by custom patterns.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// For runtime errors, we don't want to show usage
		cmd.SilenceUsage = true
		color.NoColor = false // Force color output

		if logsPattern == "" {
			return fmt.Errorf("pattern is required")
		}

		grepper, err := log.NewLogGrepper()
		if err != nil {
			return fmt.Errorf("failed to create log grepper: %v", err)
		}

		var namespaces []string
		if logsNamespace != "" {
			// Split, trim and skip empty tokens; deduplicate while preserving order
			seen := make(map[string]struct{})
			for _, tok := range strings.Split(logsNamespace, ",") {
				t := strings.TrimSpace(tok)
				if t == "" {
					continue
				}
				if _, ok := seen[t]; ok {
					continue
				}
				seen[t] = struct{}{}
				namespaces = append(namespaces, t)
			}
			if len(namespaces) == 0 {
				return fmt.Errorf("invalid namespace list: no valid namespaces provided")
			}
		} else {
			namespaces = []string{""}
		}

		var messages []log.Message
		for _, ns := range namespaces {
			ns = strings.TrimSpace(ns)
			var nsMessages []log.Message
			if ns != "" {
				if logsResource != "" {
					nsMessages, err = grepper.Grep(ns, logsResource, logsPattern, logsSortBy)
				} else {
					nsMessages, err = grepper.GrepNamespace(ns, logsPattern, logsSortBy)
				}
			} else {
				if logsResource != "" {
					nsMessages, err = grepper.GrepResourceWithoutNamespace(logsResource, logsPattern, logsSortBy)
				} else {
					nsMessages, err = grepper.GrepWithoutNamespace(logsPattern, logsSortBy)
				}
			}

			if err != nil {
				if ns != "" {
					return fmt.Errorf("failed to search logs in namespace %q: %v", ns, err)
				}
			}
			messages = append(messages, nsMessages...)
		}

		if len(namespaces) > 1 {
			messages = grepper.SortMessages(messages, logsSortBy)
		}

		printLogMessages(messages, logsPattern)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().StringVarP(&logsNamespace, "namespace", "n", "", "The Kubernetes namespace(s), comma-separated")
	logsCmd.Flags().StringVarP(&logsResource, "resource", "r", "", "The Kubernetes resource name")
	logsCmd.Flags().StringVarP(&logsPattern, "pattern", "p", "", "grep search pattern")
	logsCmd.Flags().StringVarP(&logsSortBy, "sort-by", "s", "timestamp", "Sort by: timestamp, message")

	if err := logsCmd.MarkFlagRequired("pattern"); err != nil {
		panic(fmt.Sprintf("failed to mark pattern flag as required: %v", err))
	}
}

func printLogMessages(messages []log.Message, pattern string) {
	if len(messages) == 0 {
		return
	}

	boldRed := color.New(color.FgRed).Add(color.Bold)

	for _, message := range messages {
		highlightedMessage := strings.ReplaceAll(message.Message, pattern, boldRed.Sprint(pattern))
		prefix := color.BlueString("%s/%s[%d]:", message.PodName, message.ContainerName, message.LineNumber)
		fmt.Printf("%s %s\n", prefix, highlightedMessage)
	}
}
