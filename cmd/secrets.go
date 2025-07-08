package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hbelmiro/kgrep/internal/resource"
	"github.com/spf13/cobra"
)

var (
	secretsNamespace     string
	secretsPattern       string
	secretsAllNamespaces bool
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Search Secrets in Kubernetes",
	Long:  `Search the content of Secrets for specific patterns within designated namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// For runtime errors, we don't want to show usage
		cmd.SilenceUsage = true
		color.NoColor = false // Force color output

		if secretsPattern == "" {
			return fmt.Errorf("pattern is required")
		}

		if secretsAllNamespaces && secretsNamespace != "" {
			return fmt.Errorf("--all-namespaces and --namespace cannot be used together")
		}

		resourceSearcher, err := resource.NewResourceSearcher("secrets")
		if err != nil {
			return fmt.Errorf("failed to create resource searcher: %v", err)
		}

		var occurrences []resource.Occurrence
		if secretsAllNamespaces {
			occurrences, err = resourceSearcher.SearchAllNamespaces(secretsPattern)
			if err != nil {
				return fmt.Errorf("failed to search secrets: %v", err)
			}
		} else if secretsNamespace != "" {
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
	secretsCmd.Flags().BoolVarP(&secretsAllNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces")

	if err := secretsCmd.MarkFlagRequired("pattern"); err != nil {
		panic(fmt.Sprintf("failed to mark pattern flag as required: %v", err))
	}
}
