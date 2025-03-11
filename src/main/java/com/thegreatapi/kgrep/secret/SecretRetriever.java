package com.thegreatapi.kgrep.secret;

import com.thegreatapi.kgrep.resource.ResourceRetriever;
import io.fabric8.kubernetes.api.model.Secret;
import io.fabric8.kubernetes.client.KubernetesClient;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import java.util.List;

@ApplicationScoped
final class SecretRetriever extends ResourceRetriever<Secret> {

    private final KubernetesClient client;

    @Inject
    SecretRetriever(KubernetesClient client) {
        this.client = client;
    }

    @Override
    public List<Secret> getResources(String namespace) {
        return client.secrets().inNamespace(namespace).list().getItems();
    }
}
