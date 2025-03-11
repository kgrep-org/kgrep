package com.thegreatapi.kgrep.secret;

import com.thegreatapi.kgrep.VersionProvider;
import com.thegreatapi.kgrep.resource.AbstractResourceCommand;
import io.fabric8.kubernetes.api.model.Secret;
import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "secrets", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class SecretsCommand extends AbstractResourceCommand<Secret> implements Runnable {

    @Inject
    SecretsCommand(SecretRetriever retriever, SecretGrepper grepper) {
        super(retriever, grepper);
    }
}
