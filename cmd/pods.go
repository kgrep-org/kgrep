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
	Long:  `Search the content of Kubernetes Pods for specific patterns within designated namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if podsPattern == "" {
			return fmt.Errorf("pattern is required")
		}

		resourceSearcher := resource.NewResourceSearcher("pods")
		var occurrences []resource.Occurrence

		if podsNamespace != "" {
			occurrences = resourceSearcher.Search(podsNamespace, podsPattern)
		} else {
			occurrences = resourceSearcher.SearchWithoutNamespace(podsPattern)
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
