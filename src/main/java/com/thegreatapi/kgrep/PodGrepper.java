package com.thegreatapi.kgrep;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.fabric8.kubernetes.api.model.Pod;
import io.fabric8.kubernetes.api.model.PodList;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.fabric8.kubernetes.client.dsl.MixedOperation;
import io.fabric8.kubernetes.client.dsl.PodResource;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

@ApplicationScoped
class PodGrepper extends ResourceGrepper<Pod, PodList, PodResource> {

    private final KubernetesClient client;

    @Inject
    PodGrepper(KubernetesClient client, ObjectMapper mapper, Grep grep) {
        super(mapper, grep);
        this.client = client;
    }

    @Override
    MixedOperation<Pod, PodList, PodResource> getResources() {
        return client.pods();
    }
}
