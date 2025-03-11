package com.thegreatapi.kgrep.secret;

import com.thegreatapi.kgrep.KubernetesTestsUtil;
import com.thegreatapi.kgrep.resource.ResourceLine;
import io.fabric8.kubernetes.api.model.Secret;
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
class SecretGrepperTest {

    private static final String NAMESPACE = "hbelmiro";

    private static final String PATTERN = "hbelmiro";

    @Inject
    KubernetesClient client;

    @Inject
    SecretGrepper grepper;

    @Test
    void grep() {
        createSecret();

        await().atMost(20, TimeUnit.SECONDS)
                .until(() -> grepper.grep(NAMESPACE, PATTERN).size() == 2);

        assertThat(grepper.grep(NAMESPACE, PATTERN))
                .containsExactlyInAnyOrder(
                        new ResourceLine("secrets/ds-pipeline-db-dspa", 9, "      :\\\"v2\\\"},\\\"name\\\":\\\"ds-pipeline-db-dspa\\\",\\\"namespace\\\":\\\"hbelmiro\\\",\\\"ownerReferences\\\"\\"),
                        new ResourceLine("secrets/ds-pipeline-db-dspa", 21, "  namespace: \"hbelmiro\"")
                        );
    }

    private void createSecret() {
        InputStream stream = KubernetesTestsUtil.getResourceAsStream("secret.yaml");
        Secret secret = client.secrets().load(stream).item();

        client.secrets().resource(secret).create();
    }
}