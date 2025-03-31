package com.thegreatapi.kgrep.resource;

import com.thegreatapi.kgrep.VersionProvider;
import io.fabric8.kubernetes.api.model.GenericKubernetesResource;
import jakarta.inject.Inject;
import picocli.CommandLine;

import java.util.List;

@CommandLine.Command(name = "resources", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class ResourcesCommand extends AbstractResourceCommand implements Runnable {

    @CommandLine.Option(names = {"-av", "--api-version"}, description = "apiVersion", defaultValue = "v1")
    private String apiVersion;

    @CommandLine.Option(names = {"-k", "--kind"}, description = "kind", required = true)
    private String kind;

    @Inject
    ResourcesCommand(GenericResourceRetriever resourceRetriever, GenericResourceGrepper resourceGrepper) {
        setResourceGrepper(resourceGrepper);
        setApiVersion(this.apiVersion);
        setResourceRetriever(resourceRetriever);
    }

    @Override
    public void run() {
        getOccurrences(namespace, pattern, apiVersion, kind).forEach(this::print);
    }

    @Override
    public List<ResourceLine> getOccurrences(String pattern, String apiVersion, String kind) {
        List<GenericKubernetesResource> resources = getResourceRetriever().getResources(apiVersion, kind);
        return getResourceGrepper().grep(kind, resources, pattern);
    }

    public List<ResourceLine> getOccurrences(String namespace, String pattern, String apiVersion, String kind) {
        List<GenericKubernetesResource> resources = getResourceRetriever().getResources(namespace, apiVersion, kind);
        return getResourceGrepper().grep(kind, resources, pattern);
    }
}
