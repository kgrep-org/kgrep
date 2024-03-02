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

@ApplicationScoped
public final class LogGrepper {

    private final KubernetesClient kubernetesClient;

    private final LogReader logReader;

    private final Grep grep;

    @Inject
    public LogGrepper(KubernetesClient kubernetesClient, LogReader logReader, Grep grep) {
        this.kubernetesClient = kubernetesClient;
        this.logReader = logReader;
        this.grep = grep;
    }

    List<LogMessage> grep(String namespace, String resource, String pattern) {
        List<LogMessage> lines = new ArrayList<>();

        kubernetesClient.pods().list().getItems().stream()
                .filter(pod -> pod.getMetadata().getName().contains(resource))
                .forEach(pod -> pod.getStatus().getContainerStatuses().stream()
                        .filter(status -> status.getState().getTerminated() == null)
                        .forEach(container -> lines.addAll(readLog(namespace, pod, container, pattern))));

        return Collections.unmodifiableList(lines);
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
