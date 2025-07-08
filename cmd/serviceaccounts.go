package cmd

import (
	"fmt"

	"github.com/hbelmiro/kgrep/internal/resource"
	"github.com/spf13/cobra"
)

var (
	serviceaccountsNamespace     string
	serviceaccountsPattern       string
	serviceaccountsAllNamespaces bool
)

var serviceaccountsCmd = &cobra.Command{
	Use:   "serviceaccounts",
	Short: "Search ServiceAccounts in Kubernetes",
	Long:  `Search the content of ServiceAccounts for specific patterns within designated namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// For runtime errors, we don't want to show usage
		cmd.SilenceUsage = true
		if serviceaccountsPattern == "" {
			return fmt.Errorf("pattern is required")
		}

		if serviceaccountsAllNamespaces && serviceaccountsNamespace != "" {
			return fmt.Errorf("--all-namespaces and --namespace cannot be used together")
		}

		resourceSearcher, err := resource.NewResourceSearcher("serviceaccounts")
		if err != nil {
			return fmt.Errorf("failed to create resource searcher: %v", err)
		}

		var occurrences []resource.Occurrence
		if serviceaccountsAllNamespaces {
			occurrences, err = resourceSearcher.SearchAllNamespaces(serviceaccountsPattern)
			if err != nil {
				return fmt.Errorf("failed to search serviceaccounts: %v", err)
			}
		} else if serviceaccountsNamespace != "" {
			occurrences, err = resourceSearcher.Search(serviceaccountsNamespace, serviceaccountsPattern)
			if err != nil {
				return fmt.Errorf("failed to search serviceaccounts: %v", err)
			}
		} else {
			occurrences, err = resourceSearcher.SearchWithoutNamespace(serviceaccountsPattern)
			if err != nil {
				return fmt.Errorf("failed to search serviceaccounts: %v", err)
			}
		}

		printResourceOccurrences(occurrences, serviceaccountsPattern)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serviceaccountsCmd)

	serviceaccountsCmd.Flags().StringVarP(&serviceaccountsNamespace, "namespace", "n", "", "The Kubernetes namespace")
	serviceaccountsCmd.Flags().StringVarP(&serviceaccountsPattern, "pattern", "p", "", "grep search pattern")
	serviceaccountsCmd.Flags().BoolVarP(&serviceaccountsAllNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces")

	if err := serviceaccountsCmd.MarkFlagRequired("pattern"); err != nil {
		panic(fmt.Sprintf("failed to mark pattern flag as required: %v", err))
	}
}
