package com.thegreatapi.kgrep;

import io.fabric8.kubernetes.client.KubernetesClient;
import jakarta.enterprise.inject.Produces;
import jakarta.inject.Inject;

class TestModeKgrepProducer {

    @Inject
    KubernetesClient kubernetesClient;

    @Inject
    Grep grep;

    @Inject
    @TestMode
    FakeLogReader fakeLogReader;

    @Produces
    @TestMode
    LogGrepper produceKgrep() {
        return new LogGrepper(kubernetesClient, fakeLogReader, grep);
    }
}
