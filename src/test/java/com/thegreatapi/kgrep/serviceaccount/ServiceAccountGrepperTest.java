package com.thegreatapi.kgrep.serviceaccount;

import com.thegreatapi.kgrep.KubernetesTestsUtil;
import com.thegreatapi.kgrep.resource.ResourceLine;
import io.fabric8.kubernetes.api.model.ServiceAccount;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.quarkus.test.junit.QuarkusTest;
import io.quarkus.test.kubernetes.client.WithKubernetesTestServer;
import jakarta.inject.Inject;
import org.junit.jupiter.api.Test;

import java.io.InputStream;
import java.util.concurrent.TimeUnit;

import static org.assertj.core.api.Assertions.assertThat;
import static org.awaitility.Awaitility.await;

@WithKubernetesTestServer
@QuarkusTest
class ServiceAccountGrepperTest {

    private static final String NAMESPACE = "kubeflow";

    @Inject
    KubernetesClient client;

    @Inject
    ServiceAccountGrepper serviceAccountGrepper;

    @Test
    void grep() {
        createServiceAccount();

        await().atMost(20, TimeUnit.SECONDS)
                .until(() -> serviceAccountGrepper.grep(NAMESPACE, "kubeflow").size() == 3);

        assertThat(serviceAccountGrepper.grep(NAMESPACE, "kubeflow"))
                .containsExactlyInAnyOrder(
                        new ResourceLine("serviceaccounts/pipeline-runner", 8, "      :\\\"kubeflow-pipelines\\\"},\\\"name\\\":\\\"pipeline-runner\\\",\\\"namespace\\\":\\\"kubeflow\\\"\\"),
                        new ResourceLine("serviceaccounts/pipeline-runner", 13, "    application-crd-id: \"kubeflow-pipelines\""),
                        new ResourceLine("serviceaccounts/pipeline-runner", 15, "  namespace: \"kubeflow\"")
                );
    }

    private void createServiceAccount() {
        InputStream stream = KubernetesTestsUtil.getResourceAsStream("serviceaccount.yaml");
        ServiceAccount serviceAccount = client.serviceAccounts().load(stream).item();

        client.serviceAccounts().resource(serviceAccount).create();
    }
}