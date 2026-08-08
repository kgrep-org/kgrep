package resource

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
)

// wrapKubernetesError takes an error returned from a Kubernetes API call and returns a more user-friendly error message if it recognizes the error type If err is nil, it returns nil. If the error is not recognized, it returns the original error.
func wrapKubernetesError(err error) error {
	if err == nil {
		return nil
	}

	if errors.IsUnauthorized(err) {
		return fmt.Errorf("you are not authorized to access the cluster. Please check if you are logged in and your credentials are valid: %w", err)
	}

	if errors.IsForbidden(err) {
		return fmt.Errorf("you do not have permission to perform this action in the cluster: %w", err)
	}

	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "no configuration has been provided") ||
		strings.Contains(errStr, "unable to load in-cluster configuration") ||
		strings.Contains(errStr, "couldn't find kubeconfig file") {
		return fmt.Errorf("no Kubernetes configuration found. Please ensure you have a valid kubeconfig file or that your KUBECONFIG environment variable is set: %w", err)
	}

	return err
}
