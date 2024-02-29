package com.thegreatapi.kgrep;

import com.fasterxml.jackson.core.JsonProcessingException;
import picocli.CommandLine;

import java.text.MessageFormat;
import java.util.List;

abstract class AbstractResourceCommand implements Runnable {

    private static final String ANSI_RESET = "\u001B[0m";

    private static final String BLUE = "\033[0;34m";

    @CommandLine.Option(names = {"-n", "--namespace"}, description = "The Kubernetes namespace", required = true)
    protected String namespace;

    @CommandLine.Option(names = {"-p", "--pattern"}, description = "grep search pattern", required = true)
    protected String pattern;

    private final ResourceGrepper<?, ?, ?> resourceGrepper;

    AbstractResourceCommand(ResourceGrepper<?, ?, ?> resourceGrepper) {
        this.resourceGrepper = resourceGrepper;
    }

    @Override
    public void run() {
        try {
            List<ResourceLine> occurrences = resourceGrepper.grep(namespace, pattern);
            occurrences.forEach(AbstractResourceCommand::print);
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new RuntimeException(e);
        }
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
