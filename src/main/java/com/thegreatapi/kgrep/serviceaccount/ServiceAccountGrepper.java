package com.thegreatapi.kgrep.serviceaccount;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.thegreatapi.kgrep.grep.Grep;
import com.thegreatapi.kgrep.resource.AbstractResourceGrepper;
import io.fabric8.kubernetes.api.model.ServiceAccount;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

@ApplicationScoped
final class ServiceAccountGrepper extends AbstractResourceGrepper<ServiceAccount> {

    @SuppressWarnings("unused")
    // Needed for CDI
    private ServiceAccountGrepper() {
        super(null, null);
    }

    @Inject
    ServiceAccountGrepper(ObjectMapper mapper, Grep grep) {
        super(mapper, grep);
    }
}
