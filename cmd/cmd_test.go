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

	// Check that the output contains the expected text
	if !bytes.Contains([]byte(output), []byte("kgrep is a command-line utility designed to simplify the process of searching and analyzing logs and resources in Kubernetes")) {
		t.Errorf("Expected output to contain help text, got: %s", output)
	}
}

func TestResourcesCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "resources")
	if err == nil || !strings.Contains(err.Error(), "required flag(s)") {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}

func TestResourcesCommand_WithFlags(t *testing.T) {
	output, err := executeCommand(rootCmd, "resources", "--pattern", "test", "--kind", "Pod")
	if err != nil {
		// This is expected to fail since we don't have a real kubectl connection
		// But we can check that the command structure is correct
		if !strings.Contains(err.Error(), "Error getting default namespace") &&
			!strings.Contains(err.Error(), "Error getting resources") &&
			!strings.Contains(err.Error(), "auto-discovery failed") {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	_ = output // Use output to avoid unused variable warning
}

func TestPodsCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "pods")
	if err == nil || err.Error() != "required flag(s) \"pattern\" not set" {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}

func TestPodsCommand_WithFlags(t *testing.T) {
	output, err := executeCommand(rootCmd, "pods", "--pattern", "test")
	if err != nil {
		// This is expected to fail since we don't have a real kubectl connection
		// But we can check that the command structure is correct
		if !strings.Contains(err.Error(), "Error getting default namespace") &&
			!strings.Contains(err.Error(), "Error getting resources") {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	_ = output // Use output to avoid unused variable warning
}

func TestConfigMapsCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "configmaps")
	if err == nil || err.Error() != "required flag(s) \"pattern\" not set" {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}

func TestConfigMapsCommand_WithFlags(t *testing.T) {
	output, err := executeCommand(rootCmd, "configmaps", "--pattern", "test")
	if err != nil {
		// This is expected to fail since we don't have a real kubectl connection
		// But we can check that the command structure is correct
		if !strings.Contains(err.Error(), "Error getting default namespace") &&
			!strings.Contains(err.Error(), "Error getting resources") {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	_ = output // Use output to avoid unused variable warning
}

func TestSecretsCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "secrets")
	if err == nil || err.Error() != "required flag(s) \"pattern\" not set" {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}

func TestSecretsCommand_WithFlags(t *testing.T) {
	output, err := executeCommand(rootCmd, "secrets", "--pattern", "test")
	if err != nil {
		// This is expected to fail since we don't have a real kubectl connection
		// But we can check that the command structure is correct
		if !strings.Contains(err.Error(), "Error getting default namespace") &&
			!strings.Contains(err.Error(), "Error getting resources") {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	_ = output // Use output to avoid unused variable warning
}

func TestServiceAccountsCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "serviceaccounts")
	if err == nil || err.Error() != "required flag(s) \"pattern\" not set" {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}

func TestServiceAccountsCommand_WithFlags(t *testing.T) {
	output, err := executeCommand(rootCmd, "serviceaccounts", "--pattern", "test")
	if err != nil {
		// This is expected to fail since we don't have a real kubectl connection
		// But we can check that the command structure is correct
		if !strings.Contains(err.Error(), "Error getting default namespace") &&
			!strings.Contains(err.Error(), "Error getting resources") {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	_ = output // Use output to avoid unused variable warning
}

func TestLogsCommand_MissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "logs")
	if err == nil || err.Error() != "required flag(s) \"pattern\" not set" {
		t.Errorf("Expected error for missing required flags, got: %v", err)
	}
}

func TestLogsCommand_WithFlags(t *testing.T) {
	output, err := executeCommand(rootCmd, "logs", "--pattern", "test")
	if err != nil {
		// This is expected to fail since we don't have a real kubectl connection
		// But we can check that the command structure is correct
		if !strings.Contains(err.Error(), "Error getting default namespace") &&
			!strings.Contains(err.Error(), "Error getting pods") {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	_ = output // Use output to avoid unused variable warning
}

func TestLogsCommand_WithNamespaceAndResource(t *testing.T) {
	output, err := executeCommand(rootCmd, "logs", "--pattern", "test", "--namespace", "default", "--resource", "my-pod")
	if err != nil {
		// This is expected to fail since we don't have a real kubectl connection
		// But we can check that the command structure is correct
		if !strings.Contains(err.Error(), "Error getting pods") {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	_ = output // Use output to avoid unused variable warning
}

func TestLogsCommand_WithSortBy(t *testing.T) {
	output, err := executeCommand(rootCmd, "logs", "--pattern", "test", "--sort-by", "MESSAGE")
	if err != nil {
		// This is expected to fail since we don't have a real kubectl connection
		// But we can check that the command structure is correct
		if !strings.Contains(err.Error(), "Error getting default namespace") &&
			!strings.Contains(err.Error(), "Error getting pods") {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	_ = output // Use output to avoid unused variable warning
}
