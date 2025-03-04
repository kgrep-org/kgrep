package com.thegreatapi.kgrep;

import com.thegreatapi.kgrep.configmap.ConfigMapsCommand;
import com.thegreatapi.kgrep.log.LogsCommand;
import com.thegreatapi.kgrep.pod.PodsCommand;
import com.thegreatapi.kgrep.serviceaccount.ServiceAccountsCommand;
import io.quarkus.picocli.runtime.annotations.TopCommand;
import picocli.CommandLine.Command;

@TopCommand
@Command(name = "kgrep", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class, subcommands = {
        LogsCommand.class,
        ConfigMapsCommand.class,
        PodsCommand.class,
        ServiceAccountsCommand.class
})
class KgrepCommand {
}
