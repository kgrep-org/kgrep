package com.thegreatapi.kgrep.resource;

import io.fabric8.kubernetes.api.model.GenericKubernetesResource;
import picocli.CommandLine;

import java.text.MessageFormat;
import java.util.List;

public abstract class AbstractResourceCommand implements Runnable {

    private static final String ANSI_RESET = "\u001B[0m";

    private static final String BLUE = "\033[0;34m";

    @CommandLine.Option(names = {"-n", "--namespace"}, description = "The Kubernetes namespace", required = true)
    protected String namespace;

    @CommandLine.Option(names = {"-p", "--pattern"}, description = "grep search pattern", required = true)
    protected String pattern;

    private GenericResourceRetriever resourceRetriever;

    private GenericResourceGrepper resourceGrepper;

    private String kind;

    private String apiVersion;

    protected AbstractResourceCommand() {
    }

    @Override
    public void run() {
        getOccurrences(this.kind, this.namespace, this.pattern).forEach(this::print);
    }

    public List<ResourceLine> getOccurrences(String kind, String namespace, String pattern) {
        List<GenericKubernetesResource> resources = this.resourceRetriever.getResources(namespace, apiVersion, kind);
        return this.resourceGrepper.grep(kind, resources, pattern);
    }

    protected final void print(ResourceLine resourceLine) {
        String output = MessageFormat.format("{0}{1}[{2}]:{3} {4}",
                BLUE,
                resourceLine.resourceName(),
                resourceLine.lineNumber(),
                ANSI_RESET,
                resourceLine.text()
        );

        System.out.println(output);
    }

    protected final GenericResourceGrepper getResourceGrepper() {
        return resourceGrepper;
    }

    protected final void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    protected final void setKind(String kind) {
        this.kind = kind;
    }

    protected final void setResourceGrepper(GenericResourceGrepper resourceGrepper) {
        this.resourceGrepper = resourceGrepper;
    }

    protected final void setResourceRetriever(GenericResourceRetriever resourceRetriever) {
        this.resourceRetriever = resourceRetriever;
    }

    protected final GenericResourceRetriever getResourceRetriever() {
        return resourceRetriever;
    }
}
