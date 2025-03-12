package com.thegreatapi.kgrep.configmap;

import com.thegreatapi.kgrep.resource.ResourceRetriever;
import io.fabric8.kubernetes.api.model.ConfigMap;
import io.fabric8.kubernetes.client.KubernetesClient;
import jakarta.enterprise.context.ApplicationScoped;

import java.util.List;

@ApplicationScoped
final class ConfigMapRetriever extends ResourceRetriever<ConfigMap> {

    private final KubernetesClient client;

    ConfigMapRetriever(KubernetesClient client) {
        this.client = client;
    }

    @Override
    public List<ConfigMap> getResources(String namespace) {
        return client.configMaps().inNamespace(namespace).list().getItems();
    }
}
