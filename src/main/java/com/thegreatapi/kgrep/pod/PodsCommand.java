package com.thegreatapi.kgrep.pod;

import com.thegreatapi.kgrep.VersionProvider;
import com.thegreatapi.kgrep.resource.AbstractResourceCommand;
import io.fabric8.kubernetes.api.model.Pod;
import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "pods", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class PodsCommand extends AbstractResourceCommand<Pod> implements Runnable {

    @Inject
    PodsCommand(PodRetriever retriever, PodGrepper grepper) {
        super(retriever, grepper);
    }
}
