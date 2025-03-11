package com.thegreatapi.kgrep.configmap;

import com.thegreatapi.kgrep.VersionProvider;
import com.thegreatapi.kgrep.resource.AbstractResourceCommand;
import io.fabric8.kubernetes.api.model.ConfigMap;
import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "configmaps", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class ConfigMapsCommand extends AbstractResourceCommand<ConfigMap> implements Runnable {

    @Inject
    ConfigMapsCommand(ConfigMapRetriever retriever, ConfigMapGrepper grepper) {
        super(retriever, grepper);
    }
}
