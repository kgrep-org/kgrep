package com.thegreatapi.kgrep;

import io.fabric8.kubernetes.client.KubernetesClient;
import jakarta.enterprise.inject.Produces;
import jakarta.inject.Inject;

class TestModeKgrepProducer {

    @Inject
    KubernetesClient kubernetesClient;

    @Inject
    @TestMode
    FakeLogReader fakeLogReader;

    @Produces
    @TestMode
    Kgrep produceKgrep() {
        return new Kgrep(kubernetesClient, fakeLogReader);
    }
}
