package com.thegreatapi.kgrep.configmap;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.thegreatapi.kgrep.grep.Grep;
import com.thegreatapi.kgrep.resource.AbstractResourceGrepper;
import io.fabric8.kubernetes.api.model.ConfigMap;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

@ApplicationScoped
final class ConfigMapGrepper extends AbstractResourceGrepper<ConfigMap> {

    @SuppressWarnings("unused")
    // Needed for CDI
    private ConfigMapGrepper() {
        super(null, null);
    }

    @Inject
    ConfigMapGrepper(ObjectMapper mapper, Grep grep) {
        super(mapper, grep);
    }
}
