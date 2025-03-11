package com.thegreatapi.kgrep.resource;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.thegreatapi.kgrep.grep.Grep;
import com.thegreatapi.kgrep.grep.Occurrence;
import io.fabric8.kubernetes.api.model.HasMetadata;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public abstract class AbstractResourceGrepper<T extends HasMetadata> {

    private final ObjectMapper mapper;

    private final Grep grep;

    protected AbstractResourceGrepper(ObjectMapper mapper, Grep grep) {
        this.mapper = mapper;
        this.grep = grep;
    }

    public List<ResourceLine> grep(List<T> resources, String pattern) {
        List<ResourceLine> occurrences = new ArrayList<>();

        try (ExecutorService executorService = Executors.newVirtualThreadPerTaskExecutor()) {
            for (T resource : resources) {
                executorService.execute(() -> {
                    String[] lines = getYaml(resource).split(System.lineSeparator());

                    grep.run(lines, pattern).stream()
                            .map(occurrence -> createResourceLine(resource, occurrence))
                            .forEach(occurrences::add);
                });
            }
        }

        return occurrences;
    }

    private String getYaml(T resource) {
        try {
            return mapper.writeValueAsString(resource);
        } catch (JsonProcessingException e) {
            throw new RuntimeException("Error while getting " + resource.getMetadata().getName() + " as YAML.", e);
        }
    }

    private ResourceLine createResourceLine(T resource, Occurrence occurrence) {
        return new ResourceLine(resource.getFullResourceName() + "/" + resource.getMetadata().getName(),
                occurrence.lineNumber(), occurrence.text());
    }
}
