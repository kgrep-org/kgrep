package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hbelmiro/kgrep/internal/resource"
	"github.com/spf13/cobra"
)

var (
	secretsNamespace string
	secretsPattern   string
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Search Secrets in Kubernetes",
	Long:  `Search the content of Secrets for specific patterns within designated namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		color.NoColor = false // Force color output

		if secretsPattern == "" {
			return fmt.Errorf("pattern is required")
		}

		resourceSearcher, err := resource.NewResourceSearcher("secrets")
		if err != nil {
			return fmt.Errorf("failed to create resource searcher: %v", err)
		}

		var occurrences []resource.Occurrence
		if secretsNamespace != "" {
			occurrences, err = resourceSearcher.Search(secretsNamespace, secretsPattern)
			if err != nil {
				return fmt.Errorf("failed to search secrets: %v", err)
			}
		} else {
			occurrences, err = resourceSearcher.SearchWithoutNamespace(secretsPattern)
			if err != nil {
				return fmt.Errorf("failed to search secrets: %v", err)
			}
		}

		printResourceOccurrences(occurrences, secretsPattern)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(secretsCmd)

	secretsCmd.Flags().StringVarP(&secretsNamespace, "namespace", "n", "", "The Kubernetes namespace")
	secretsCmd.Flags().StringVarP(&secretsPattern, "pattern", "p", "", "grep search pattern")

	if err := secretsCmd.MarkFlagRequired("pattern"); err != nil {
		panic(fmt.Sprintf("failed to mark pattern flag as required: %v", err))
	}
}
