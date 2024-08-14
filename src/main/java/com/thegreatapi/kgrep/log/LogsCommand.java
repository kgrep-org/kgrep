package com.thegreatapi.kgrep.log;

import com.thegreatapi.kgrep.VersionProvider;
import jakarta.inject.Inject;
import picocli.CommandLine;

import java.text.MessageFormat;
import java.util.Optional;

@CommandLine.Command(name = "logs", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class LogsCommand implements Runnable {

    private static final String ANSI_RESET = "\u001B[0m";

    private static final String PURPLE = "\033[0;35m";

    private static final String BLUE = "\033[0;34m";

    @CommandLine.Option(names = {"-n", "--namespace"}, description = "The Kubernetes namespace", required = true)
    private String namespace;

    @CommandLine.Option(names = {"-r", "--resource-name"}, description = "Resource")
    private Optional<String> resource;

    @CommandLine.Option(names = {"-p", "--pattern"}, description = "grep search pattern", required = true)
    private String pattern;

    @CommandLine.Option(names = {"-s", "--sort-by"}, description = "sort by", defaultValue = "POD_AND_CONTAINER")
    private SortBy sortBy;

    private final LogGrepper logGrepper;

    @SuppressWarnings("unused")
    @Inject
    LogsCommand(LogGrepper logGrepper) {
        this.logGrepper = logGrepper;
    }

    @Override
    public void run() {
        if (resource.isPresent()) {
            logGrepper.grep(namespace, resource.orElseThrow(), pattern, sortBy)
                    .forEach(LogsCommand::print);
        } else {
            logGrepper.grep(namespace, pattern, sortBy)
                    .forEach(LogsCommand::print);
        }
    }

    private static void print(LogMessage logMessage) {
        String output = MessageFormat.format("{0}{1}{2}/{3}[{4}]:{5} {6}",
                BLUE,
                logMessage.podName(),
                PURPLE,
                logMessage.containerName(),
                logMessage.lineNumber(),
                ANSI_RESET,
                logMessage.message());

        System.out.println(output);
    }
}
