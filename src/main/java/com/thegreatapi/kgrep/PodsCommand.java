package com.thegreatapi.kgrep;

import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "pods", mixinStandardHelpOptions = true)
class PodsCommand extends AbstractResourceCommand implements Runnable {

    @Inject
    PodsCommand(PodGrepper podGrepper) {
        super(podGrepper);
    }
}
