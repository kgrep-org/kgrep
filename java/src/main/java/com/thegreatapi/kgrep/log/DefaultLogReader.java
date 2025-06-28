package com.thegreatapi.kgrep.log;

import io.fabric8.kubernetes.client.KubernetesClient;
import jakarta.enterprise.context.ApplicationScoped;

@ApplicationScoped
class DefaultLogReader implements LogReader {

    private final KubernetesClient kubernetesClient;

    DefaultLogReader(KubernetesClient kubernetesClient) {
        this.kubernetesClient = kubernetesClient;
    }

    @Override
    public String read(String namespace, String podName, String containerName) {
        return kubernetesClient.pods().inNamespace(namespace).withName(podName).inContainer(containerName).getLog();
    }
}
