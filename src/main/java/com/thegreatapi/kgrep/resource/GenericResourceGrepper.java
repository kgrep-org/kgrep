package com.thegreatapi.kgrep.resource;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.thegreatapi.kgrep.grep.Grep;
import com.thegreatapi.kgrep.grep.Occurrence;
import io.fabric8.kubernetes.api.model.GenericKubernetesResource;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

@ApplicationScoped
public final class GenericResourceGrepper {

    private final ObjectMapper mapper;

    private final Grep grep;

    @Inject
    GenericResourceGrepper(ObjectMapper mapper, Grep grep) {
        this.mapper = mapper;
        this.grep = grep;
    }

    public List<ResourceLine> grep(String kind, List<GenericKubernetesResource> resources, String pattern) {
        List<ResourceLine> occurrences = new ArrayList<>();

        try (ExecutorService executorService = Executors.newVirtualThreadPerTaskExecutor()) {
            for (GenericKubernetesResource resource : resources) {
                executorService.execute(() -> {
                    String[] lines = getYaml(resource).split(System.lineSeparator());

                    grep.run(lines, pattern).stream()
                            .map(occurrence -> createResourceLine(kind, resource, occurrence))
                            .forEach(occurrences::add);
                });
            }
        }

        return occurrences;
    }

    private String getYaml(GenericKubernetesResource resource) {
        try {
            return mapper.writeValueAsString(resource);
        } catch (JsonProcessingException e) {
            throw new RuntimeException("Error while getting " + resource.getMetadata().getName() + " as YAML.", e);
        }
    }

    private ResourceLine createResourceLine(String kind, GenericKubernetesResource resource, Occurrence occurrence) {
        return new ResourceLine(kind + "/" + resource.getMetadata().getName(),
                occurrence.lineNumber(), occurrence.text());
    }
}
