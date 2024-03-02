package com.thegreatapi.kgrep;

import picocli.CommandLine;

import java.io.IOException;
import java.io.InputStream;
import java.util.Objects;

public final class VersionProvider implements CommandLine.IVersionProvider {

    private final String[] version;

    public VersionProvider() throws IOException {
        try (InputStream stream = VersionProvider.class.getResourceAsStream("/version")) {
            this.version = new String[]{new String(Objects.requireNonNull(stream).readAllBytes())};
        }
    }

    @Override
    public String[] getVersion() {
        return version;
    }
}
