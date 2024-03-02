package com.thegreatapi.kgrep.pod;

import com.thegreatapi.kgrep.VersionProvider;
import com.thegreatapi.kgrep.resource.AbstractResourceCommand;
import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "pods", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class PodsCommand extends AbstractResourceCommand implements Runnable {

    @Inject
    PodsCommand(PodGrepper podGrepper) {
        super(podGrepper);
    }
}
