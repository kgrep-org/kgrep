package com.thegreatapi.kgrep.resource;

import io.fabric8.kubernetes.api.model.APIResource;
import io.fabric8.kubernetes.api.model.APIResourceList;
import io.fabric8.kubernetes.api.model.GenericKubernetesResource;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.fabric8.kubernetes.client.dsl.base.ResourceDefinitionContext;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import java.util.List;
import java.util.Objects;
import java.util.Optional;

@ApplicationScoped
public final class GenericResourceRetriever {

    private final KubernetesClient client;

    @Inject
    GenericResourceRetriever(KubernetesClient client) {
        this.client = client;
    }

    public List<GenericKubernetesResource> getResources(String apiVersion, String kind) {
        Objects.requireNonNull(apiVersion);
        Objects.requireNonNull(kind);

        String namespace = client.getNamespace();
        if (namespace == null) {
            namespace = client.getConfiguration().getNamespace();
        }
        return getResources(namespace, apiVersion, kind);
    }

    public List<GenericKubernetesResource> getResources(String namespace, String apiVersion, String kind) {
        Objects.requireNonNull(namespace);
        Objects.requireNonNull(apiVersion);
        Objects.requireNonNull(kind);

        APIResourceList resourceList = client.getApiResources(apiVersion);

        Optional<APIResource> foundResource = resourceList.getResources().stream()
                .filter(resource -> resource.getKind().equalsIgnoreCase(kind))
                .findFirst();

        if (foundResource.isPresent()) {
            String correctKind = foundResource.get().getKind();
            return getResourceDefinitionContext(namespace, apiVersion, correctKind);
        } else {
            return getResourceDefinitionContext(namespace, apiVersion, kind);
        }
    }

    private List<GenericKubernetesResource> getResourceDefinitionContext(String namespace, String apiVersion, String kind) {
        ResourceDefinitionContext context = new ResourceDefinitionContext.Builder()
                .withGroup(apiVersion.contains("/") ? apiVersion.split("/")[0] : "")
                .withVersion(apiVersion.contains("/") ? apiVersion.split("/")[1] : apiVersion)
                .withKind(kind)
                .withNamespaced(true)
                .build();

        return client.genericKubernetesResources(context).inNamespace(namespace).list().getItems();
    }
}
