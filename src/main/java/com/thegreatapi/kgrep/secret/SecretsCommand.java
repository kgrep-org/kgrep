package com.thegreatapi.kgrep.secret;

import com.thegreatapi.kgrep.VersionProvider;
import com.thegreatapi.kgrep.resource.AbstractResourceCommand;
import com.thegreatapi.kgrep.resource.GenericResourceGrepper;
import com.thegreatapi.kgrep.resource.GenericResourceRetriever;
import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "secrets", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class SecretsCommand extends AbstractResourceCommand implements Runnable {

    @Inject
    SecretsCommand(GenericResourceRetriever retriever, GenericResourceGrepper grepper) {
        setResourceRetriever(retriever);
        setResourceGrepper(grepper);
        setApiVersion("v1");
        setKind("Secret");
    }
}
