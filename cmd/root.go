package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kgrep",
	Short:   "kgrep - Search and analyze logs and resources in Kubernetes",
	Long:    `kgrep is a command-line utility designed to simplify the process of searching and analyzing logs and resources in Kubernetes.`,
	Version: "0.1.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
