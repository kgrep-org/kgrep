package com.thegreatapi.kgrep;

import io.fabric8.kubernetes.api.model.ContainerStatus;
import io.fabric8.kubernetes.api.model.Pod;
import io.fabric8.kubernetes.client.KubernetesClient;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.stream.Stream;

@ApplicationScoped
class Kgrep {

    private final KubernetesClient kubernetesClient;

    private final LogReader logReader;

    @Inject
    Kgrep(KubernetesClient kubernetesClient, LogReader logReader) {
        this.kubernetesClient = kubernetesClient;
        this.logReader = logReader;
    }

    List<LogMessage> run(String namespace, String resource, String grep) {
//            List<ConfigMap> configMaps = client.configMaps()/*.inNamespace(namespace)*/.list().getItems();
//            System.out.println("ConfigMaps in Namespace " + namespace + ":");
//            for (ConfigMap configMap : configMaps) {
//                System.out.println(configMap.getMetadata().getName());
//            }

        List<LogMessage> lines = new ArrayList<>();

        kubernetesClient.pods().list().getItems().stream()
                .filter(pod -> pod.getMetadata().getName().contains(resource))
                .forEach(pod -> pod.getStatus().getContainerStatuses().stream()
                        .filter(status -> status.getState().getTerminated() == null)
                        .forEach(container -> lines.addAll(readLog(namespace, pod, container, grep))));

        return Collections.unmodifiableList(lines);
    }

    private List<LogMessage> readLog(String namespace, Pod pod, ContainerStatus container, String grep) {
        String podName = pod.getMetadata().getName();

        String log = logReader.read(namespace, podName, container.getName());

        return Stream.of(log.split(System.lineSeparator()))
                .filter(line -> line.contains(grep))
                .map(line -> new LogMessage(podName, container.getName(), line))
                .toList();
    }
}
