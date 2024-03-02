package com.thegreatapi.kgrep.log;

import com.thegreatapi.kgrep.TestMode;
import jakarta.enterprise.context.ApplicationScoped;

import java.util.HashMap;
import java.util.Map;

@ApplicationScoped
@TestMode
public final class FakeLogReader implements LogReader {

    private final Map<RegistryKey, String> logRegistry = new HashMap<>();


    @Override
    public String read(String namespace, String podName, String containerName) {
        return logRegistry.getOrDefault(new RegistryKey(namespace, podName, containerName), "");
    }

    void addLog(RegistryKey registryKey, String log) {
        this.logRegistry.put(registryKey, log);
    }

    record RegistryKey(String namespace, String podName, String containerName) {
    }
}
