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
		color.NoColor = false // Force color output

		if logsPattern == "" {
			return fmt.Errorf("pattern is required")
		}

		grepper, err := log.NewLogGrepper()
		if err != nil {
			return fmt.Errorf("failed to create log grepper: %v", err)
		}

		var messages []log.Message

		if logsNamespace != "" {
			if logsResource != "" {
				messages, err = grepper.Grep(logsNamespace, logsResource, logsPattern, logsSortBy)
				if err != nil {
					return fmt.Errorf("failed to search logs: %v", err)
				}
			} else {
				messages, err = grepper.GrepNamespace(logsNamespace, logsPattern, logsSortBy)
				if err != nil {
					return fmt.Errorf("failed to search logs: %v", err)
				}
			}
		} else {
			if logsResource != "" {
				messages, err = grepper.GrepResourceWithoutNamespace(logsResource, logsPattern, logsSortBy)
				if err != nil {
					return fmt.Errorf("failed to search logs: %v", err)
				}
			} else {
				messages, err = grepper.GrepWithoutNamespace(logsPattern, logsSortBy)
				if err != nil {
					return fmt.Errorf("failed to search logs: %v", err)
				}
			}
		}

		printLogMessages(messages, logsPattern)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().StringVarP(&logsNamespace, "namespace", "n", "", "The Kubernetes namespace")
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

	blue := "\033[34m"
	reset := "\033[0m"
	boldRed := "\033[1;31m"

	for _, message := range messages {
		highlightedMessage := strings.ReplaceAll(message.Message, pattern, boldRed+pattern+reset)
		prefix := fmt.Sprintf("%s%s/%s[%d]:%s", blue, message.PodName, message.ContainerName, message.LineNumber, reset)
		fmt.Printf("%s %s\n", prefix, highlightedMessage)
	}
}
