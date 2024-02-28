package com.thegreatapi.kgrep;

import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "logs", mixinStandardHelpOptions = true)
class LogsCommand implements Runnable {

    private static final String ANSI_RESET = "\u001B[0m";

    private static final String PURPLE = "\033[0;35m";

    private static final String BLUE = "\033[0;34m";

    @CommandLine.Option(names = {"-n", "--namespace"}, description = "The Kubernetes namespace", required = true)
    private String namespace;

    @CommandLine.Option(names = {"-r", "--resource-name"}, description = "Resource", required = true)
    private String resource;

    @CommandLine.Option(names = {"-p", "--pattern"}, description = "grep search pattern", required = true)
    private String pattern;

    private final LogGrepper logGrepper;

    @SuppressWarnings("unused")
    @Inject
    LogsCommand(LogGrepper logGrepper) {
        this.logGrepper = logGrepper;
    }

    @Override
    public void run() {
        logGrepper.grep(namespace, resource, pattern)
                .forEach(LogsCommand::print);
    }

    private static void print(LogMessage logMessage) {
        System.out.println(BLUE + logMessage.podName() + PURPLE + "/" + logMessage.containerName() + ANSI_RESET
                + ": " + logMessage.message());
    }
}
