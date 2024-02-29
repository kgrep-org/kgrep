package com.thegreatapi.kgrep;

import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "configmaps", mixinStandardHelpOptions = true)
class ConfigMapsCommand extends AbstractResourceCommand implements Runnable {

    @Inject
    ConfigMapsCommand(ConfigMapGrepper configMapGrepper) {
        super(configMapGrepper);
    }
}
