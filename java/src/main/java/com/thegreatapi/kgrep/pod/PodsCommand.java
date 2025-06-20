package com.thegreatapi.kgrep.pod;

import com.thegreatapi.kgrep.VersionProvider;
import com.thegreatapi.kgrep.resource.AbstractResourceCommand;
import com.thegreatapi.kgrep.resource.GenericResourceGrepper;
import com.thegreatapi.kgrep.resource.GenericResourceRetriever;
import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "pods", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class PodsCommand extends AbstractResourceCommand implements Runnable {

    @Inject
    PodsCommand(GenericResourceRetriever retriever, GenericResourceGrepper grepper) {
        setResourceRetriever(retriever);
        setResourceGrepper(grepper);
        setApiVersion("v1");
        setKind("Pod");
    }
}
