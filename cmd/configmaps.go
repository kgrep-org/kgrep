package cmd

import (
	"fmt"

	"github.com/hbelmiro/kgrep/internal/resource"
	"github.com/spf13/cobra"
)

var (
	configmapsNamespace string
	configmapsPattern   string
)

var configmapsCmd = &cobra.Command{
	Use:   "configmaps",
	Short: "Search ConfigMaps in Kubernetes",
	Long:  `Search the content of ConfigMaps for specific patterns within designated namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if configmapsPattern == "" {
			return fmt.Errorf("pattern is required")
		}

		resourceSearcher := resource.NewResourceSearcher("configmaps")
		var occurrences []resource.Occurrence

		if configmapsNamespace != "" {
			occurrences = resourceSearcher.Search(configmapsNamespace, configmapsPattern)
		} else {
			occurrences = resourceSearcher.SearchWithoutNamespace(configmapsPattern)
		}

		printResourceOccurrences(occurrences, configmapsPattern)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configmapsCmd)

	configmapsCmd.Flags().StringVarP(&configmapsNamespace, "namespace", "n", "", "The Kubernetes namespace")
	configmapsCmd.Flags().StringVarP(&configmapsPattern, "pattern", "p", "", "grep search pattern")

	if err := configmapsCmd.MarkFlagRequired("pattern"); err != nil {
		panic(fmt.Sprintf("failed to mark pattern flag as required: %v", err))
	}
}
