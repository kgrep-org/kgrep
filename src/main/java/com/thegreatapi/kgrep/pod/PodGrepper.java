package com.thegreatapi.kgrep.pod;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.thegreatapi.kgrep.grep.Grep;
import com.thegreatapi.kgrep.resource.AbstractResourceGrepper;
import io.fabric8.kubernetes.api.model.Pod;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

@ApplicationScoped
final class PodGrepper extends AbstractResourceGrepper<Pod> {

    @SuppressWarnings("unused")
    // Needed for CDI
    private PodGrepper() {
        super(null, null);
    }

    @Inject
    PodGrepper(ObjectMapper mapper, Grep grep) {
        super(mapper, grep);
    }
}
