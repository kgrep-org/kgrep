package com.thegreatapi.kgrep.grep;

import com.thegreatapi.kgrep.TestMode;
import com.thegreatapi.kgrep.grep.Grep;
import com.thegreatapi.kgrep.log.FakeLogReader;
import com.thegreatapi.kgrep.log.LogGrepper;
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
