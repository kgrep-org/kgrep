package com.thegreatapi.kgrep.serviceaccount;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.thegreatapi.kgrep.grep.Grep;
import com.thegreatapi.kgrep.resource.ResourceGrepper;
import io.fabric8.kubernetes.api.model.ServiceAccount;
import io.fabric8.kubernetes.api.model.ServiceAccountList;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.fabric8.kubernetes.client.dsl.MixedOperation;
import io.fabric8.kubernetes.client.dsl.ServiceAccountResource;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

@ApplicationScoped
public class ServiceAccountGrepper extends ResourceGrepper<ServiceAccount, ServiceAccountList, ServiceAccountResource> {

    private final KubernetesClient client;

    @Inject
    ServiceAccountGrepper(KubernetesClient client, ObjectMapper mapper, Grep grep) {
        super(mapper, grep);
        this.client = client;
    }

    @Override
    public MixedOperation<ServiceAccount, ServiceAccountList, ServiceAccountResource> getResources() {
        return client.serviceAccounts();
    }
}
