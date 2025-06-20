package cmd

import (
	"fmt"

	"github.com/hbelmiro/kgrep/internal/resource"
	"github.com/spf13/cobra"
)

var (
	serviceAccountsNamespace string
	serviceAccountsPattern   string
)

var serviceAccountsCmd = &cobra.Command{
	Use:   "serviceaccounts",
	Short: "Search ServiceAccounts in Kubernetes",
	Long:  `Search the content of Kubernetes ServiceAccounts for specific patterns within designated namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if serviceAccountsPattern == "" {
			return fmt.Errorf("pattern is required")
		}

		resourceSearcher := resource.NewResourceSearcher("serviceaccounts")
		var occurrences []resource.Occurrence

		if serviceAccountsNamespace != "" {
			occurrences = resourceSearcher.Search(serviceAccountsNamespace, serviceAccountsPattern)
		} else {
			occurrences = resourceSearcher.SearchWithoutNamespace(serviceAccountsPattern)
		}

		for _, occ := range occurrences {
			fmt.Printf("serviceaccounts/%s[%d]: %s\n", occ.Resource, occ.Line, occ.Content)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serviceAccountsCmd)

	serviceAccountsCmd.Flags().StringVarP(&serviceAccountsNamespace, "namespace", "n", "", "The Kubernetes namespace")
	serviceAccountsCmd.Flags().StringVarP(&serviceAccountsPattern, "pattern", "p", "", "grep search pattern")

	if err := serviceAccountsCmd.MarkFlagRequired("pattern"); err != nil {
		panic(fmt.Sprintf("failed to mark pattern flag as required: %v", err))
	}
}
