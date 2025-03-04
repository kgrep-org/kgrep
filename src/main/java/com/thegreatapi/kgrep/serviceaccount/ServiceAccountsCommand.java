package com.thegreatapi.kgrep.serviceaccount;

import com.thegreatapi.kgrep.VersionProvider;
import com.thegreatapi.kgrep.resource.AbstractResourceCommand;
import jakarta.inject.Inject;
import picocli.CommandLine;

@CommandLine.Command(name = "serviceaccounts", mixinStandardHelpOptions = true, versionProvider = VersionProvider.class)
public final class ServiceAccountsCommand extends AbstractResourceCommand implements Runnable {

    @Inject
    ServiceAccountsCommand(ServiceAccountGrepper serviceAccountGrepper) {
        super(serviceAccountGrepper);
    }
}
