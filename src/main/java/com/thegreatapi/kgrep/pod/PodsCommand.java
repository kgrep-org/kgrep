package com.thegreatapi.kgrep.pod;

import com.thegreatapi.kgrep.resource.AbstractResourceCommand;
import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "pods", mixinStandardHelpOptions = true)
public final class PodsCommand extends AbstractResourceCommand implements Runnable {

    @Inject
    PodsCommand(PodGrepper podGrepper) {
        super(podGrepper);
    }
}
