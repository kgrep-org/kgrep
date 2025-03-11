package com.thegreatapi.kgrep.secret;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.thegreatapi.kgrep.grep.Grep;
import com.thegreatapi.kgrep.resource.AbstractResourceGrepper;
import io.fabric8.kubernetes.api.model.Secret;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

@ApplicationScoped
final class SecretGrepper extends AbstractResourceGrepper<Secret> {

    @SuppressWarnings("unused")
    // Needed for CDI
    private SecretGrepper() {
        super(null, null);
    }

    @Inject
    SecretGrepper(ObjectMapper mapper, Grep grep) {
        super(mapper, grep);
    }
}
