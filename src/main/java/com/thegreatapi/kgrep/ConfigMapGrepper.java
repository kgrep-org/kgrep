package com.thegreatapi.kgrep;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.fabric8.kubernetes.api.model.ConfigMap;
import io.fabric8.kubernetes.client.KubernetesClient;
import jakarta.enterprise.context.ApplicationScoped;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

@ApplicationScoped
class ConfigMapGrepper {

    private final KubernetesClient kubernetesClient;

    private final ObjectMapper mapper;

    private final Grep grep;

    ConfigMapGrepper(KubernetesClient kubernetesClient, ObjectMapper mapper, Grep grep) {
        this.kubernetesClient = kubernetesClient;
        this.mapper = mapper;
        this.grep = grep;
    }

    List<ResourceLine> grep(String namespace, String pattern) throws JsonProcessingException, InterruptedException {
        List<ConfigMap> configMaps = kubernetesClient.configMaps().inNamespace(namespace).list().getItems();

        List<ResourceLine> occurrences = new ArrayList<>();

        try (ExecutorService executorService = Executors.newVirtualThreadPerTaskExecutor()) {
            for (ConfigMap configMap : configMaps) {
                executorService.submit(() -> {
                    String[] lines = getYaml(configMap).split(System.lineSeparator());

                    grep.run(lines, pattern).stream()
                            .map(line -> createResourceLine(configMap, line))
                            .forEach(occurrences::add);
                });
            }

            executorService.shutdown();

            if (!executorService.awaitTermination(1, TimeUnit.MINUTES)) {
                throw new RuntimeException("Timeout!!");
            }
        }

        return occurrences;
    }

    private String getYaml(ConfigMap configMap) {
        try {
            return mapper.writeValueAsString(configMap);
        } catch (JsonProcessingException e) {
            throw new RuntimeException("Error while getting " + configMap.getMetadata().getName() + " as YAML.", e);
        }
    }

    private static ResourceLine createResourceLine(ConfigMap configMap, String line) {
        return new ResourceLine(configMap.getFullResourceName() + "/" + configMap.getMetadata().getName(),
                line);
    }
}
