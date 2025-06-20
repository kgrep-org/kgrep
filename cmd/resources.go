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
	Short: "Search any Kubernetes resource",
	Long:  `Search the content of any Kubernetes resource for specific patterns within designated namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resourcesPattern == "" {
			return fmt.Errorf("pattern is required")
		}

		if resourcesKind == "" {
			return fmt.Errorf("kind is required")
		}

		var resourceSearcher *resource.Searcher

		if resourcesAPIVersion != "" {
			// Use explicit API version and kind
			resourceSearcher = resource.NewGenericResourceSearcher(resourcesAPIVersion, resourcesKind)
		} else {
			// Use auto-discovery for API version
			resourceSearcher = resource.NewAutoDiscoveryResourceSearcher(resourcesKind)
		}

		var occurrences []resource.Occurrence

		if resourcesNamespace != "" {
			occurrences = resourceSearcher.Search(resourcesNamespace, resourcesPattern)
		} else {
			occurrences = resourceSearcher.SearchWithoutNamespace(resourcesPattern)
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
