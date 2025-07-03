package cmd

import (
	"fmt"

	"github.com/hbelmiro/kgrep/internal/resource"
	"github.com/spf13/cobra"
)

var (
	podsNamespace string
	podsPattern   string
)

var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Search Pods in Kubernetes",
	Long:  `Search the content of Pods for specific patterns within designated namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if podsPattern == "" {
			return fmt.Errorf("pattern is required")
		}

		resourceSearcher, err := resource.NewResourceSearcher("pods")
		if err != nil {
			return fmt.Errorf("failed to create resource searcher: %v", err)
		}

		var occurrences []resource.Occurrence
		if podsNamespace != "" {
			occurrences, err = resourceSearcher.Search(podsNamespace, podsPattern)
			if err != nil {
				return fmt.Errorf("failed to search pods: %v", err)
			}
		} else {
			occurrences, err = resourceSearcher.SearchWithoutNamespace(podsPattern)
			if err != nil {
				return fmt.Errorf("failed to search pods: %v", err)
			}
		}

		printResourceOccurrences(occurrences, podsPattern)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(podsCmd)

	podsCmd.Flags().StringVarP(&podsNamespace, "namespace", "n", "", "The Kubernetes namespace")
	podsCmd.Flags().StringVarP(&podsPattern, "pattern", "p", "", "grep search pattern")

	if err := podsCmd.MarkFlagRequired("pattern"); err != nil {
		panic(fmt.Sprintf("failed to mark pattern flag as required: %v", err))
	}
}
