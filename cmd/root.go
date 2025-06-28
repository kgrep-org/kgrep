package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// These variables are set at build time using ldflags
var (
	Version    = "dev"
	BuildTime  = "unknown"
	CommitHash = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "kgrep",
	Short:   "kgrep - Search and analyze logs and resources in Kubernetes",
	Long:    `kgrep is a command-line utility designed to simplify the process of searching and analyzing logs and resources in Kubernetes.`,
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add a custom version command that shows more details
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show detailed version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("kgrep version %s\n", Version)
			fmt.Printf("Build time: %s\n", BuildTime)
			fmt.Printf("Commit hash: %s\n", CommitHash)
		},
	}
	rootCmd.AddCommand(versionCmd)
}
