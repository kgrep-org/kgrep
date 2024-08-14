package com.thegreatapi.kgrep.log;

import com.thegreatapi.kgrep.grep.Grep;
import com.thegreatapi.kgrep.grep.Occurrence;
import io.fabric8.kubernetes.api.model.ContainerStatus;
import io.fabric8.kubernetes.api.model.Pod;
import io.fabric8.kubernetes.client.KubernetesClient;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.function.Predicate;

import static java.util.Comparator.comparing;

@ApplicationScoped
public final class LogGrepper {

    private static final Predicate<Pod> ALL_PODS = _ -> true;

    private final KubernetesClient kubernetesClient;

    private final LogReader logReader;

    private final Grep grep;

    @Inject
    public LogGrepper(KubernetesClient kubernetesClient, LogReader logReader, Grep grep) {
        this.kubernetesClient = kubernetesClient;
        this.logReader = logReader;
        this.grep = grep;
    }

    List<LogMessage> grep(String namespace, String pattern) {
        return grep(namespace, pattern, SortBy.POD_AND_CONTAINER);
    }

    List<LogMessage> grep(String namespace, String resource, String pattern) {
        return grep(namespace, resource, pattern, SortBy.POD_AND_CONTAINER);
    }

    List<LogMessage> grep(String namespace, String pattern, SortBy sortBy) {
        return grep(namespace, pattern, sortBy, ALL_PODS);
    }

    List<LogMessage> grep(String namespace, String resource, String pattern, SortBy sortBy) {
        return grep(namespace, pattern, sortBy, pod -> pod.getMetadata().getName().contains(resource));
    }

    private List<LogMessage> grep(String namespace, String pattern, SortBy sortBy, Predicate<Pod> podFilter) {
        List<LogMessage> lines = new ArrayList<>();

        for (Pod pod : kubernetesClient.pods().inNamespace(namespace).list().getItems()) {
            if (podFilter.test(pod)) {
                for (ContainerStatus status : pod.getStatus().getContainerStatuses()) {
                    if (status.getState().getWaiting() == null) {
                        lines.addAll(readLog(namespace, pod, status, pattern));
                    }
                }
            }
        }

        switch (sortBy) {
            case POD_AND_CONTAINER -> {
                return Collections.unmodifiableList(lines);
            }
            case MESSAGE -> {
                lines.sort(comparing(LogMessage::message));
                return Collections.unmodifiableList(lines);
            }
            case null, default -> throw new UnsupportedOperationException("Sorting by " + sortBy + " is not supported");
        }
    }

    private List<LogMessage> readLog(String namespace, Pod pod, ContainerStatus container, String pattern) {
        String podName = pod.getMetadata().getName();

        String log = logReader.read(namespace, podName, container.getName());

        List<Occurrence> occurrences = grep.run(log.split(System.lineSeparator()), pattern);

        return occurrences.stream()
                .map(occurrence -> new LogMessage(podName, container.getName(), occurrence.text(), occurrence.lineNumber()))
                .toList();
    }
}
