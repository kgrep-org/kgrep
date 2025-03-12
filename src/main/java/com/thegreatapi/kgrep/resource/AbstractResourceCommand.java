package com.thegreatapi.kgrep.resource;

import io.fabric8.kubernetes.api.model.HasMetadata;
import picocli.CommandLine;

import java.text.MessageFormat;
import java.util.List;

public abstract class AbstractResourceCommand<T extends HasMetadata> implements Runnable {

    private static final String ANSI_RESET = "\u001B[0m";

    private static final String BLUE = "\033[0;34m";

    @CommandLine.Option(names = {"-n", "--namespace"}, description = "The Kubernetes namespace", required = true)
    protected String namespace;

    @CommandLine.Option(names = {"-p", "--pattern"}, description = "grep search pattern", required = true)
    protected String pattern;

    private final ResourceRetriever<T> resourceRetriever;

    private final AbstractResourceGrepper<T> resourceGrepper;

    protected AbstractResourceCommand(ResourceRetriever<T> resourceRetriever, AbstractResourceGrepper<T> resourceGrepper) {
        this.resourceRetriever = resourceRetriever;
        this.resourceGrepper = resourceGrepper;
    }

    @Override
    public void run() {
        getOccurrences(this.namespace, this.pattern).forEach(AbstractResourceCommand::print);
    }

    public List<ResourceLine> getOccurrences(String namespace, String pattern) {
        List<T> resources = this.resourceRetriever.getResources(namespace);
        return this.resourceGrepper.grep(resources, pattern);
    }

    private static void print(ResourceLine resourceLine) {
        String output = MessageFormat.format("{0}{1}[{2}]:{3} {4}",
                BLUE,
                resourceLine.resourceName(),
                resourceLine.lineNumber(),
                ANSI_RESET,
                resourceLine.text()
        );

        System.out.println(output);
    }
}
