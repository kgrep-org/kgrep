package com.thegreatapi.kgrep.serviceaccount;

import com.thegreatapi.kgrep.resource.ResourceRetriever;
import io.fabric8.kubernetes.api.model.ServiceAccount;
import io.fabric8.kubernetes.client.KubernetesClient;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import java.util.List;

@ApplicationScoped
final class ServiceAccountRetriever extends ResourceRetriever<ServiceAccount> {

    private final KubernetesClient client;

    @Inject
    ServiceAccountRetriever(KubernetesClient client) {
        this.client = client;
    }

    @Override
    public List<ServiceAccount> getResources(String namespace) {
        return client.serviceAccounts().inNamespace(namespace).list().getItems();
    }
}
