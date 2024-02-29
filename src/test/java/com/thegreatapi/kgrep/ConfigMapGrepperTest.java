package com.thegreatapi.kgrep;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.fabric8.kubernetes.api.model.ConfigMap;
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
class ConfigMapGrepperTest {

    private static final String NAMESPACE = "kubeflow";
    @Inject
    KubernetesClient client;

    @Inject
    ObjectMapper mapper;

    @Inject
    Grep grep;

    @Test
    void grep() throws JsonProcessingException, InterruptedException {
        createConfigMaps();

        var configMapGrepper = new ConfigMapGrepper(client, mapper, grep);

        List<ResourceLine> occurrences = configMapGrepper.grep(NAMESPACE, "kubeflow");

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
        InputStream stream = getResourceAsStream(yamlOrJson);
        ConfigMap configMap = client.configMaps().load(stream).item();

        client.configMaps().resource(configMap).create();
    }

    private static InputStream getResourceAsStream(String yamlOrJson) {
        return ConfigMapGrepperTest.class.getClassLoader().getResourceAsStream(yamlOrJson);
    }
}