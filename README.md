### kgrep

#### Overview

This utility is designed to simplify the process of searching and analyzing logs from multiple Kubernetes pods. It
takes two parameters: the `grep_search_parameter` for filtering logs and the `pod_name_keyword` to identify specific
pods.

#### Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed and configured to connect to your
  Kubernetes cluster.

#### Usage

```bash
./kgrep <grep_search_parameter> <pod_name_keyword>
```

- `<grep_search_parameter>`: The search parameter used with `grep` to filter the logs.
- `<pod_name_keyword>`: The keyword used to identify relevant pods.

#### Example

```bash
./kgrep ERROR my-app
```

#### Output

For each matching pod, the script retrieves and displays the pod name, then prints the logs containing the
specified `grep_search_parameter`. The output is separated by dashed lines for clarity.

```plaintext
Pod: my-app-pod-1
[log entry matching "ERROR" parameter]
----------------------------------------
Pod: my-app-pod-2
[log entry matching "ERROR" parameter]
----------------------------------------
```

#### Notes

- Ensure that the script has executable permissions (`chmod +x kgrep`).
- This script assumes a typical `kubectl` configuration and may need adjustments based on your environment.
- Adjustments may be required for specific log formatting in your Kubernetes environment.