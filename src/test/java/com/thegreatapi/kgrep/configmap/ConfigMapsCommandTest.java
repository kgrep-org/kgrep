package com.thegreatapi.kgrep.configmap;

import com.thegreatapi.kgrep.KubernetesTestsUtil;
import com.thegreatapi.kgrep.resource.ResourceLine;
import io.fabric8.kubernetes.api.model.ConfigMap;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.quarkus.test.junit.QuarkusTest;
import io.quarkus.test.kubernetes.client.WithKubernetesTestServer;
import jakarta.inject.Inject;
import org.junit.jupiter.api.Test;

import java.io.InputStream;
import java.util.List;
import java.util.concurrent.TimeUnit;

import static org.assertj.core.api.Assertions.assertThat;
import static org.awaitility.Awaitility.await;

@WithKubernetesTestServer
@QuarkusTest
class ConfigMapsCommandTest {

    private static final String NAMESPACE = "kubeflow";

    private final KubernetesClient client;

    private final ConfigMapsCommand command;

    @Inject
    ConfigMapsCommandTest(KubernetesClient client, ConfigMapGrepper grepper, ConfigMapRetriever retriever) {
        this.command = new ConfigMapsCommand(retriever, grepper);
        this.client = client;
    }

    @Test
    void grep() {
        createConfigMaps();

        List<ResourceLine> occurrences = command.getOccurrences(NAMESPACE, "kubeflow");

        await().atMost(20, TimeUnit.SECONDS)
                .until(() -> command.getOccurrences(NAMESPACE, "kubeflow").size() == 5);

        assertThat(occurrences)
                .containsExactlyInAnyOrder(
                        new ResourceLine("configmaps/custom-ui-configmap", 8, "  namespace: \"kubeflow\""),
                        new ResourceLine("configmaps/kfp-launcher", 8, "      \\ http://minio-sample.kubeflow.svc.cluster.local:9000\\\\n  region: minio\\\\n \\"),
                        new ResourceLine("configmaps/kfp-launcher", 12, "      :\\\"data-science-pipelines\\\"},\\\"name\\\":\\\"kfp-launcher\\\",\\\"namespace\\\":\\\"kubeflow\\\"\\"),
                        new ResourceLine("configmaps/kfp-launcher", 23, "  namespace: \"kubeflow\""),
                        new ResourceLine("configmaps/kfp-launcher", 35, "  providers: \"s3:\\n  endpoint: http://minio-sample.kubeflow.svc.cluster.local:9000\\n\\")
                );
    }

    private void createConfigMaps() {
        createConfigMap("configmap1.yaml");
        createConfigMap("configmap2.yaml");
    }

    private void createConfigMap(String yamlOrJson) {
        InputStream stream = KubernetesTestsUtil.getResourceAsStream(yamlOrJson);
        ConfigMap configMap = client.configMaps().load(stream).item();

        client.configMaps().resource(configMap).create();
    }
}