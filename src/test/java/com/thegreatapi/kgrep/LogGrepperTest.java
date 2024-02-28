package com.thegreatapi.kgrep;

import com.thegreatapi.kgrep.FakeLogReader.RegistryKey;
import io.fabric8.kubernetes.api.model.ContainerState;
import io.fabric8.kubernetes.api.model.ContainerStateBuilder;
import io.fabric8.kubernetes.api.model.ContainerStatus;
import io.fabric8.kubernetes.api.model.ContainerStatusBuilder;
import io.fabric8.kubernetes.api.model.Pod;
import io.fabric8.kubernetes.api.model.PodBuilder;
import io.fabric8.kubernetes.api.model.PodStatus;
import io.fabric8.kubernetes.api.model.PodStatusBuilder;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.quarkus.test.junit.QuarkusTest;
import io.quarkus.test.kubernetes.client.WithKubernetesTestServer;
import jakarta.inject.Inject;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

@WithKubernetesTestServer
@QuarkusTest
class LogGrepperTest {

    private static final String NAMESPACE = "test";

    @Inject
    KubernetesClient client;

    @Inject
    @TestMode
    LogGrepper logGrepper;

    @Inject
    @TestMode
    FakeLogReader fakeLogReader;

    @Test
    void testInteractionWithAPIServer() {
        createPods();

        assertThat(logGrepper.grep(NAMESPACE, "pod", "initialized")).containsExactlyInAnyOrder(
                new LogMessage("pod1", "container1", "xpto initialized"),
                new LogMessage("pod2", "container2", "foo initialized"),
                new LogMessage("pod2", "container2", "bar initialized")
        );
    }

    private void createPods() {
        client.pods().resource(createPod("container1", "pod1")).create();
        client.pods().resource(createPod("container2", "pod2")).create();


        fakeLogReader.addLog(new RegistryKey(NAMESPACE, "pod1", "container1"),
                """
                        Initializing xpto
                        xpto initialized
                        error writing to xpto
                        """);

        fakeLogReader.addLog(new RegistryKey(NAMESPACE, "pod2", "container2"),
                """
                        Initializing foo
                        foo initialized
                        error writing to foo
                        Initializing bar
                        bar initialized
                        error writing to bar
                        """);
    }

    private static Pod createPod(String containerName, String podName) {
        ContainerState containerState = new ContainerStateBuilder()
                .build();

        ContainerStatus containerStatus = new ContainerStatusBuilder()
                .withName(containerName)
                .withState(containerState)
                .build();

        PodStatus status = new PodStatusBuilder().withContainerStatuses(containerStatus)
                .build();

        return new PodBuilder().withNewMetadata()
                .withName(podName)
                .withNamespace(NAMESPACE)
                .and()
                .withStatus(status)
                .build();
    }
}