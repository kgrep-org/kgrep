package com.thegreatapi.kgrep;

import com.fasterxml.jackson.core.JsonProcessingException;
import io.fabric8.kubernetes.api.model.Pod;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.quarkus.test.junit.QuarkusTest;
import io.quarkus.test.kubernetes.client.WithKubernetesTestServer;
import jakarta.inject.Inject;
import org.junit.jupiter.api.Test;

import java.io.InputStream;
import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;

@WithKubernetesTestServer
@QuarkusTest
class PodGrepperTest {

    private static final String NAMESPACE = "kubeflow";

    @Inject
    KubernetesClient client;

    @Inject
    PodGrepper podGrepper;

    @Test
    void grep() throws JsonProcessingException, InterruptedException {
        createPods();

        List<ResourceLine> occurrences = podGrepper.grep(NAMESPACE, "kubeflow");

        assertThat(occurrences)
                .containsExactlyInAnyOrder(
                        new ResourceLine("pods/ds-pipeline-sample-7b59bd7cb4-szxqb", 20, "  namespace: \"kubeflow\""),
                        new ResourceLine("pods/mariadb-sample-5bd78c456f-kffct", 20, "  namespace: \"kubeflow\""),
                        new ResourceLine("pods/ds-pipeline-sample-7b59bd7cb4-szxqb", 34, "      value: \"kubeflow\""),
                        new ResourceLine("pods/ds-pipeline-sample-7b59bd7cb4-szxqb", 45, "      value: \"mariadb-sample.kubeflow.svc.cluster.local\""),
                        new ResourceLine("pods/ds-pipeline-sample-7b59bd7cb4-szxqb", 79, "      value: \"minio-sample.kubeflow.svc.cluster.local\""),
                        new ResourceLine("pods/ds-pipeline-sample-7b59bd7cb4-szxqb", 87, "      value: \"ds-pipeline-sample.kubeflow.svc.cluster.local\""),
                        new ResourceLine("pods/ds-pipeline-sample-7b59bd7cb4-szxqb", 101, "      value: \"http://minio-sample.kubeflow.svc.cluster.local:9000\""),
                        new ResourceLine("pods/ds-pipeline-sample-7b59bd7cb4-szxqb", 194, "      kubeflow\\\"}}\""),
                        new ResourceLine("pods/ds-pipeline-sample-7b59bd7cb4-szxqb", 195, "    - \"--openshift-sar={\\\"namespace\\\":\\\"kubeflow\\\",\\\"resource\\\":\\\"routes\\\",\\\"resourceName\\\"\\")
                );
    }

    private void createPods() {
        createPod("pod1.yaml");
        createPod("pod2.yaml");
    }

    private void createPod(String yamlOrJson) {
        InputStream stream = KubernetesTestsUtil.getResourceAsStream(yamlOrJson);
        Pod pod = client.pods().load(stream).item();

        client.pods().resource(pod).create();
    }
}