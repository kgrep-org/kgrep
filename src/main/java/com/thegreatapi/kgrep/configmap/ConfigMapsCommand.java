package com.thegreatapi.kgrep.configmap;

import com.thegreatapi.kgrep.VersionProvider;
import com.thegreatapi.kgrep.resource.AbstractResourceCommand;
import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "configmaps", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class ConfigMapsCommand extends AbstractResourceCommand implements Runnable {

    @Inject
    ConfigMapsCommand(ConfigMapGrepper configMapGrepper) {
        super(configMapGrepper);
    }
}
