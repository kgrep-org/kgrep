package com.thegreatapi.kgrep;

import jakarta.inject.Inject;
import picocli.CommandLine;
import picocli.CommandLine.Command;

@Command(name = "kgrep", mixinStandardHelpOptions = true)
class KgrepCommand implements Runnable {

    private static final String ANSI_RESET = "\u001B[0m";

    private static final String PURPLE = "\033[0;35m";

    private static final String BLUE = "\033[0;34m";

    @CommandLine.Option(names = {"-n", "--namespace"}, description = "The Kubernetes namespace", required = true)
    private String namespace;

    @CommandLine.Option(names = {"-r", "--resource-name"}, description = "Resource", required = true)
    private String resource;

    @CommandLine.Option(names = {"-g", "--grep-search-parameter"}, description = "grep search parameter", required = true)
    private String grep;

    private final Kgrep kgrep;

    @SuppressWarnings("unused")
    @Inject
    KgrepCommand(Kgrep kgrep) {
        this.kgrep = kgrep;
    }

    @Override
    public void run() {
        for (LogMessage logMessage : kgrep.run(namespace, resource, grep)) {
            System.out.println(BLUE + logMessage.podName() + PURPLE + "/" + logMessage.containerName() + ANSI_RESET + ": " + logMessage.message());
        }
    }
}
