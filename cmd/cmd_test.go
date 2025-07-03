package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	output, err := executeCommand(rootCmd, "--help")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !bytes.Contains([]byte(output), []byte("kgrep is a command-line utility designed to simplify the process of searching and analyzing logs and resources in Kubernetes")) {
		t.Errorf("Expected output to contain help text, got: %s", output)
	}
}

// Test that required flags are properly validated
func TestResourcesCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "resources")
	if err == nil || !strings.Contains(err.Error(), "required flag(s)") {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}

func TestPodsCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "pods")
	if err == nil || err.Error() != "required flag(s) \"pattern\" not set" {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}

func TestConfigMapsCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "configmaps")
	if err == nil || err.Error() != "required flag(s) \"pattern\" not set" {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}

func TestSecretsCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "secrets")
	if err == nil || err.Error() != "required flag(s) \"pattern\" not set" {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}

func TestServiceAccountsCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "serviceaccounts")
	if err == nil || err.Error() != "required flag(s) \"pattern\" not set" {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}

func TestLogsCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "logs")
	if err == nil || err.Error() != "required flag(s) \"pattern\" not set" {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}
