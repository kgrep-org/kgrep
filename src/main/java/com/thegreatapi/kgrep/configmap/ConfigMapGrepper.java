package com.thegreatapi.kgrep.configmap;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.thegreatapi.kgrep.grep.Grep;
import com.thegreatapi.kgrep.resource.ResourceGrepper;
import io.fabric8.kubernetes.api.model.ConfigMap;
import io.fabric8.kubernetes.api.model.ConfigMapList;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.fabric8.kubernetes.client.dsl.MixedOperation;
import io.fabric8.kubernetes.client.dsl.Resource;
import jakarta.enterprise.context.ApplicationScoped;

@ApplicationScoped
class ConfigMapGrepper extends ResourceGrepper<ConfigMap, ConfigMapList, Resource<ConfigMap>> {

    private final KubernetesClient client;

    ConfigMapGrepper(KubernetesClient client, ObjectMapper mapper, Grep grep) {
        super(mapper, grep);
        this.client = client;
    }

    @Override
    public MixedOperation<ConfigMap, ConfigMapList, Resource<ConfigMap>> getResources() {
        return client.configMaps();
    }
}
