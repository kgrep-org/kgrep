package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func resetFlags() {
	resourcesKind = ""
	resourcesNamespace = ""
	resourcesPattern = ""
	resourcesAllNamespaces = false

	podsNamespace = ""
	podsPattern = ""
	podsAllNamespaces = false

	configmapsNamespace = ""
	configmapsPattern = ""
	configmapsAllNamespaces = false

	secretsNamespace = ""
	secretsPattern = ""
	secretsAllNamespaces = false

	serviceaccountsNamespace = ""
	serviceaccountsPattern = ""
	serviceaccountsAllNamespaces = false

	logsNamespace = ""
	logsResource = ""
	logsPattern = ""
	logsSortBy = ""
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	resetFlags()
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

// Test that syntax errors show usage text
func TestResourcesCommand_MissingFlags_ShowsUsage(t *testing.T) {
	output, err := executeCommand(rootCmd, "resources")
	if err == nil {
		t.Errorf("Expected error for missing required flags")
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage text to be shown for syntax errors, got: %s", output)
	}

	if !strings.Contains(output, "kgrep resources [flags]") {
		t.Errorf("Expected command usage line to be shown, got: %s", output)
	}

	if !strings.Contains(output, "Flags:") {
		t.Errorf("Expected flags section to be shown, got: %s", output)
	}
}

func TestPodsCommand_MissingFlags_ShowsUsage(t *testing.T) {
	output, err := executeCommand(rootCmd, "pods")
	if err == nil {
		t.Errorf("Expected error for missing required flags")
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage text to be shown for syntax errors, got: %s", output)
	}

	if !strings.Contains(output, "kgrep pods [flags]") {
		t.Errorf("Expected command usage line to be shown, got: %s", output)
	}
}

func TestConfigMapsCommand_MissingFlags_ShowsUsage(t *testing.T) {
	output, err := executeCommand(rootCmd, "configmaps")
	if err == nil {
		t.Errorf("Expected error for missing required flags")
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage text to be shown for syntax errors, got: %s", output)
	}

	if !strings.Contains(output, "kgrep configmaps [flags]") {
		t.Errorf("Expected command usage line to be shown, got: %s", output)
	}
}

func TestSecretsCommand_MissingFlags_ShowsUsage(t *testing.T) {
	output, err := executeCommand(rootCmd, "secrets")
	if err == nil {
		t.Errorf("Expected error for missing required flags")
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage text to be shown for syntax errors, got: %s", output)
	}

	if !strings.Contains(output, "kgrep secrets [flags]") {
		t.Errorf("Expected command usage line to be shown, got: %s", output)
	}
}

func TestServiceAccountsCommand_MissingFlags_ShowsUsage(t *testing.T) {
	output, err := executeCommand(rootCmd, "serviceaccounts")
	if err == nil {
		t.Errorf("Expected error for missing required flags")
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage text to be shown for syntax errors, got: %s", output)
	}

	if !strings.Contains(output, "kgrep serviceaccounts [flags]") {
		t.Errorf("Expected command usage line to be shown, got: %s", output)
	}
}

func TestLogsCommand_MissingFlags_ShowsUsage(t *testing.T) {
	output, err := executeCommand(rootCmd, "logs")
	if err == nil {
		t.Errorf("Expected error for missing required flags")
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage text to be shown for syntax errors, got: %s", output)
	}

	if !strings.Contains(output, "kgrep logs [flags]") {
		t.Errorf("Expected command usage line to be shown, got: %s", output)
	}
}

// Test that runtime errors (correct syntax) don't show usage text
func TestResourcesCommand_RuntimeError_NoUsage(t *testing.T) {
	// This test simulates a runtime error with correct syntax
	// We can't easily test actual Kubernetes connectivity errors in unit tests,
	// but we can test that SilenceUsage works properly by using a command
	// that will fail during execution rather than during flag validation

	// Using a non-existent resource type will cause a runtime error
	output, err := executeCommand(rootCmd, "resources", "--kind", "NonExistentResource", "--pattern", "test")
	if err == nil {
		t.Errorf("Expected error for non-existent resource")
	}

	// The error should be present but no usage should be shown
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected error message to be shown, got: %s", output)
	}

	if strings.Contains(output, "Usage:") {
		t.Errorf("Expected no usage text for runtime errors, got: %s", output)
	}

	if strings.Contains(output, "kgrep resources [flags]") {
		t.Errorf("Expected no command usage line for runtime errors, got: %s", output)
	}

	if strings.Contains(output, "Flags:") {
		t.Errorf("Expected no flags section for runtime errors, got: %s", output)
	}
}

func TestPodsCommand_RuntimeError_NoUsage(t *testing.T) {
	output, err := executeCommand(rootCmd, "pods", "--pattern", "test", "--namespace", "non-existent-namespace")

	// Whether there's an error or not, usage should not be shown for correct syntax
	if strings.Contains(output, "Usage:") {
		t.Errorf("Expected no usage text for commands with correct syntax, got: %s", output)
	}

	if strings.Contains(output, "kgrep pods [flags]") {
		t.Errorf("Expected no command usage line for commands with correct syntax, got: %s", output)
	}

	// If there is an error, it should be a runtime error, not a syntax error
	if err != nil && strings.Contains(err.Error(), "required flag") {
		t.Errorf("Unexpected syntax error for command with correct syntax: %v", err)
	}
}

func TestConfigMapsCommand_RuntimeError_NoUsage(t *testing.T) {
	output, err := executeCommand(rootCmd, "configmaps", "--pattern", "test", "--namespace", "non-existent-namespace")

	// Whether there's an error or not, usage should not be shown for correct syntax
	if strings.Contains(output, "Usage:") {
		t.Errorf("Expected no usage text for commands with correct syntax, got: %s", output)
	}

	if strings.Contains(output, "kgrep configmaps [flags]") {
		t.Errorf("Expected no command usage line for commands with correct syntax, got: %s", output)
	}

	// If there is an error, it should be a runtime error, not a syntax error
	if err != nil && strings.Contains(err.Error(), "required flag") {
		t.Errorf("Unexpected syntax error for command with correct syntax: %v", err)
	}
}

func TestSecretsCommand_RuntimeError_NoUsage(t *testing.T) {
	output, err := executeCommand(rootCmd, "secrets", "--pattern", "test", "--namespace", "non-existent-namespace")

	// Whether there's an error or not, usage should not be shown for correct syntax
	if strings.Contains(output, "Usage:") {
		t.Errorf("Expected no usage text for commands with correct syntax, got: %s", output)
	}

	if strings.Contains(output, "kgrep secrets [flags]") {
		t.Errorf("Expected no command usage line for commands with correct syntax, got: %s", output)
	}

	// If there is an error, it should be a runtime error, not a syntax error
	if err != nil && strings.Contains(err.Error(), "required flag") {
		t.Errorf("Unexpected syntax error for command with correct syntax: %v", err)
	}
}

func TestServiceAccountsCommand_RuntimeError_NoUsage(t *testing.T) {
	output, err := executeCommand(rootCmd, "serviceaccounts", "--pattern", "test", "--namespace", "non-existent-namespace")

	// Whether there's an error or not, usage should not be shown for correct syntax
	if strings.Contains(output, "Usage:") {
		t.Errorf("Expected no usage text for commands with correct syntax, got: %s", output)
	}

	if strings.Contains(output, "kgrep serviceaccounts [flags]") {
		t.Errorf("Expected no command usage line for commands with correct syntax, got: %s", output)
	}

	// If there is an error, it should be a runtime error, not a syntax error
	if err != nil && strings.Contains(err.Error(), "required flag") {
		t.Errorf("Unexpected syntax error for command with correct syntax: %v", err)
	}
}

func TestLogsCommand_RuntimeError_NoUsage(t *testing.T) {
	output, err := executeCommand(rootCmd, "logs", "--pattern", "test", "--namespace", "non-existent-namespace")

	// Whether there's an error or not, usage should not be shown for correct syntax
	if strings.Contains(output, "Usage:") {
		t.Errorf("Expected no usage text for commands with correct syntax, got: %s", output)
	}

	if strings.Contains(output, "kgrep logs [flags]") {
		t.Errorf("Expected no command usage line for commands with correct syntax, got: %s", output)
	}

	// If there is an error, it should be a runtime error, not a syntax error
	if err != nil && strings.Contains(err.Error(), "required flag") {
		t.Errorf("Unexpected syntax error for command with correct syntax: %v", err)
	}
}

func TestAllNamespacesFlagValidation_Pods(t *testing.T) {
	output, err := executeCommand(rootCmd, "pods", "--pattern", "test", "--namespace", "test-ns", "--all-namespaces")
	if err == nil {
		t.Errorf("Expected error when using both --namespace and --all-namespaces flags")
	}

	if !strings.Contains(output, "--all-namespaces and --namespace cannot be used together") {
		t.Errorf("Expected mutual exclusion error message, got: %s", output)
	}
}

func TestAllNamespacesFlagValidation_ConfigMaps(t *testing.T) {
	output, err := executeCommand(rootCmd, "configmaps", "--pattern", "test", "--namespace", "test-ns", "--all-namespaces")
	if err == nil {
		t.Errorf("Expected error when using both --namespace and --all-namespaces flags")
	}

	if !strings.Contains(output, "--all-namespaces and --namespace cannot be used together") {
		t.Errorf("Expected mutual exclusion error message, got: %s", output)
	}
}

func TestAllNamespacesFlagValidation_Secrets(t *testing.T) {
	output, err := executeCommand(rootCmd, "secrets", "--pattern", "test", "--namespace", "test-ns", "--all-namespaces")
	if err == nil {
		t.Errorf("Expected error when using both --namespace and --all-namespaces flags")
	}

	if !strings.Contains(output, "--all-namespaces and --namespace cannot be used together") {
		t.Errorf("Expected mutual exclusion error message, got: %s", output)
	}
}

func TestAllNamespacesFlagValidation_ServiceAccounts(t *testing.T) {
	output, err := executeCommand(rootCmd, "serviceaccounts", "--pattern", "test", "--namespace", "test-ns", "--all-namespaces")
	if err == nil {
		t.Errorf("Expected error when using both --namespace and --all-namespaces flags")
	}

	if !strings.Contains(output, "--all-namespaces and --namespace cannot be used together") {
		t.Errorf("Expected mutual exclusion error message, got: %s", output)
	}
}

func TestAllNamespacesFlagValidation_Resources(t *testing.T) {
	output, err := executeCommand(rootCmd, "resources", "--kind", "Pod", "--pattern", "test", "--namespace", "test-ns", "--all-namespaces")
	if err == nil {
		t.Errorf("Expected error when using both --namespace and --all-namespaces flags")
	}

	if !strings.Contains(output, "--all-namespaces and --namespace cannot be used together") {
		t.Errorf("Expected mutual exclusion error message, got: %s", output)
	}
}

func TestAllNamespacesFlagAccepted_Pods(t *testing.T) {
	output, err := executeCommand(rootCmd, "pods", "--pattern", "test", "--all-namespaces")

	// We don't expect a flag validation error, though the command may fail for other reasons (like kubeconfig)
	if err != nil && strings.Contains(err.Error(), "--all-namespaces and --namespace cannot be used together") {
		t.Errorf("Unexpected flag validation error when using only --all-namespaces: %v, output: %s", err, output)
	}

	// Any other errors are acceptable (like kubeconfig issues)
	t.Logf("Command output: %s", output)
	if err != nil {
		t.Logf("Expected error for kubeconfig/connectivity issues: %v", err)
	}
}
