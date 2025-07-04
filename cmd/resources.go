package cmd

import (
	"fmt"

	"github.com/hbelmiro/kgrep/internal/resource"
	"github.com/spf13/cobra"
)

var (
	resourcesNamespace  string
	resourcesPattern    string
	resourcesAPIVersion string
	resourcesKind       string
)

var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Search Generic Resources in Kubernetes",
	Long:  `Search the content of any Kubernetes resource for specific patterns within designated namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// For runtime errors, we don't want to show usage
		cmd.SilenceUsage = true

		var resourceSearcher *resource.Searcher
		var err error

		if resourcesAPIVersion != "" {
			resourceSearcher, err = resource.NewGenericResourceSearcher(resourcesAPIVersion, resourcesKind)
			if err != nil {
				return fmt.Errorf("failed to create generic resource searcher: %v", err)
			}
		} else {
			resourceSearcher, err = resource.NewAutoDiscoveryResourceSearcher(resourcesKind)
			if err != nil {
				return fmt.Errorf("failed to create auto-discovery resource searcher: %v", err)
			}
		}

		var occurrences []resource.Occurrence
		if resourcesNamespace != "" {
			occurrences, err = resourceSearcher.Search(resourcesNamespace, resourcesPattern)
			if err != nil {
				return fmt.Errorf("failed to search resources: %v", err)
			}
		} else {
			occurrences, err = resourceSearcher.SearchWithoutNamespace(resourcesPattern)
			if err != nil {
				return fmt.Errorf("failed to search resources: %v", err)
			}
		}

		printResourceOccurrences(occurrences, resourcesPattern)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(resourcesCmd)

	resourcesCmd.Flags().StringVarP(&resourcesNamespace, "namespace", "n", "", "The Kubernetes namespace")
	resourcesCmd.Flags().StringVarP(&resourcesPattern, "pattern", "p", "", "grep search pattern")
	resourcesCmd.Flags().StringVar(&resourcesAPIVersion, "api-version", "", "API version (e.g., v1, apps/v1). If not provided, will be auto-discovered.")
	resourcesCmd.Flags().StringVarP(&resourcesKind, "kind", "k", "", "Resource kind (e.g., Pod, Deployment)")

	if err := resourcesCmd.MarkFlagRequired("pattern"); err != nil {
		panic(fmt.Sprintf("failed to mark pattern flag as required: %v", err))
	}
	if err := resourcesCmd.MarkFlagRequired("kind"); err != nil {
		panic(fmt.Sprintf("failed to mark kind flag as required: %v", err))
	}
}
