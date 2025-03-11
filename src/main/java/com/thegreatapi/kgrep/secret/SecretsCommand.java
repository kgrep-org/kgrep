package com.thegreatapi.kgrep.secret;

import com.thegreatapi.kgrep.VersionProvider;
import com.thegreatapi.kgrep.resource.AbstractResourceCommand;
import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "secrets", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class SecretsCommand extends AbstractResourceCommand implements Runnable {

    @Inject
    SecretsCommand(SecretGrepper grepper) {
        super(grepper);
    }
}
