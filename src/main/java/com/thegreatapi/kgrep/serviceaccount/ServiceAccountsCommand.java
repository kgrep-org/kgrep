package com.thegreatapi.kgrep.serviceaccount;

import com.thegreatapi.kgrep.VersionProvider;
import com.thegreatapi.kgrep.resource.AbstractResourceCommand;
import io.fabric8.kubernetes.api.model.ServiceAccount;
import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "serviceaccounts", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class ServiceAccountsCommand extends AbstractResourceCommand<ServiceAccount> implements Runnable {

    @Inject
    ServiceAccountsCommand(ServiceAccountRetriever serviceAccountRetriever, ServiceAccountGrepper serviceAccountGrepper) {
        super(serviceAccountRetriever, serviceAccountGrepper);
    }
}
