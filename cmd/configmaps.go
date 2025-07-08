package cmd

import (
	"fmt"

	"github.com/hbelmiro/kgrep/internal/resource"
	"github.com/spf13/cobra"
)

var (
	configmapsNamespace     string
	configmapsPattern       string
	configmapsAllNamespaces bool
)

var configmapsCmd = &cobra.Command{
	Use:   "configmaps",
	Short: "Search ConfigMaps in Kubernetes",
	Long:  `Search the content of ConfigMaps for specific patterns within designated namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// For runtime errors, we don't want to show usage
		cmd.SilenceUsage = true
		if configmapsPattern == "" {
			return fmt.Errorf("pattern is required")
		}

		if configmapsAllNamespaces && configmapsNamespace != "" {
			return fmt.Errorf("--all-namespaces and --namespace cannot be used together")
		}

		resourceSearcher, err := resource.NewResourceSearcher("configmaps")
		if err != nil {
			return fmt.Errorf("failed to create resource searcher: %v", err)
		}

		var occurrences []resource.Occurrence
		if configmapsAllNamespaces {
			occurrences, err = resourceSearcher.SearchAllNamespaces(configmapsPattern)
			if err != nil {
				return fmt.Errorf("failed to search configmaps: %v", err)
			}
		} else if configmapsNamespace != "" {
			occurrences, err = resourceSearcher.Search(configmapsNamespace, configmapsPattern)
			if err != nil {
				return fmt.Errorf("failed to search configmaps: %v", err)
			}
		} else {
			occurrences, err = resourceSearcher.SearchWithoutNamespace(configmapsPattern)
			if err != nil {
				return fmt.Errorf("failed to search configmaps: %v", err)
			}
		}

		printResourceOccurrences(occurrences, configmapsPattern)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configmapsCmd)

	configmapsCmd.Flags().StringVarP(&configmapsNamespace, "namespace", "n", "", "The Kubernetes namespace")
	configmapsCmd.Flags().StringVarP(&configmapsPattern, "pattern", "p", "", "grep search pattern")
	configmapsCmd.Flags().BoolVarP(&configmapsAllNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces")

	if err := configmapsCmd.MarkFlagRequired("pattern"); err != nil {
		panic(fmt.Sprintf("failed to mark pattern flag as required: %v", err))
	}
}
