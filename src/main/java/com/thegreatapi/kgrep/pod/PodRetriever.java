package com.thegreatapi.kgrep.pod;

import com.thegreatapi.kgrep.resource.ResourceRetriever;
import io.fabric8.kubernetes.api.model.Pod;
import io.fabric8.kubernetes.client.KubernetesClient;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import java.util.List;

@ApplicationScoped
final class PodRetriever extends ResourceRetriever<Pod> {

    private final KubernetesClient client;

    @Inject
    PodRetriever(KubernetesClient client) {
        this.client = client;
    }

    @Override
    public List<Pod> getResources(String namespace) {
        return client.pods().inNamespace(namespace).list().getItems();
    }
}
