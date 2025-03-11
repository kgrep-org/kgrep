package com.thegreatapi.kgrep.resource;

import io.fabric8.kubernetes.api.model.HasMetadata;

import java.util.List;

public abstract class ResourceRetriever<T extends HasMetadata> {

    protected ResourceRetriever() {
    }

    public abstract List<T> getResources(String namespace);

}
