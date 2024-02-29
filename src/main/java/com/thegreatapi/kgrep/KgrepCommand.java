package com.thegreatapi.kgrep;

import io.quarkus.picocli.runtime.annotations.TopCommand;
import picocli.CommandLine.Command;

@TopCommand
@Command(name = "kgrep", mixinStandardHelpOptions = true, subcommands = {
        LogsCommand.class,
        ConfigMapsCommand.class,
        PodsCommand.class
})
class KgrepCommand {
}
