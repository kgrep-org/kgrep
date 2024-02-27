### kgrep

#### Overview

This utility is designed to simplify the process of searching and analyzing logs from multiple Kubernetes pods.

#### Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed and configured to connect to your
  Kubernetes cluster.

#### Usage

```bash
kgrep [-hV] -g=<grep> -n=<namespace> -r=<resource>
  -g, --grep-search-parameter=<grep>
                  grep search parameter
  -h, --help      Show this help message and exit.
  -n, --namespace=<namespace>
                  The Kubernetes namespace
  -r, --resource-name=<resource>
                  Resource
  -V, --version   Print version information and exit.

```

#### Example

```bash
kgrep -n kubeflow -r ds-pipeline -g Success
```

#### Output

For each matching pod, the script retrieves and displays the pod name, the container name, and the log messages containing the
specified `grep-search-parameter`.

```bash
ds-pipeline-persistenceagent-sample-cbd7d67f8-tmqw7/ds-pipeline-persistenceagent: time="2024-02-27T16:02:59Z" level=info msg="Success while syncing resource (kubeflow/iris-pipeline-f0d3e)"
ds-pipeline-persistenceagent-sample-cbd7d67f8-tmqw7/ds-pipeline-persistenceagent: time="2024-02-27T16:02:59Z" level=info msg="Success while syncing resource (kubeflow/iris-pipeline-d74a4)"
ds-pipeline-persistenceagent-sample-cbd7d67f8-tmqw7/ds-pipeline-persistenceagent: time="2024-02-27T16:02:59Z" level=info msg="Success while syncing resource (kubeflow/iris-pipeline-fefd9)"
ds-pipeline-persistenceagent-sample-cbd7d67f8-tmqw7/ds-pipeline-persistenceagent: time="2024-02-27T16:02:59Z" level=info msg="Success while syncing resource (kubeflow/iris-pipeline-b8841)"
```

## Building the project

This project uses Quarkus, the Supersonic Subatomic Java Framework.

If you want to learn more about Quarkus, please visit its website: https://quarkus.io/ .

## Running the application in dev mode

You can run your application in dev mode that enables live coding using:
```shell script
./mvnw compile quarkus:dev
```

> **_NOTE:_**  Quarkus now ships with a Dev UI, which is available in dev mode only at http://localhost:8080/q/dev/.

## Packaging and running the application

The application can be packaged using:
```shell script
./mvnw package
```
It produces the `quarkus-run.jar` file in the `target/quarkus-app/` directory.
Be aware that it’s not an _über-jar_ as the dependencies are copied into the `target/quarkus-app/lib/` directory.

The application is now runnable using `java -jar target/quarkus-app/quarkus-run.jar`.

If you want to build an _über-jar_, execute the following command:
```shell script
./mvnw package -Dquarkus.package.type=uber-jar
```

The application, packaged as an _über-jar_, is now runnable using `java -jar target/*-runner.jar`.

## Creating a native executable

You can create a native executable using: 
```shell script
./mvnw package -Dnative
```

Or, if you don't have GraalVM installed, you can run the native executable build in a container using: 
```shell script
./mvnw package -Dnative -Dquarkus.native.container-build=true
```

You can then execute your native executable with: `./target/kgrep-1.0.0-SNAPSHOT-runner`

If you want to learn more about building native executables, please consult https://quarkus.io/guides/maven-tooling.