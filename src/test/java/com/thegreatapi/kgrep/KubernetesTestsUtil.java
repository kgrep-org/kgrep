package com.thegreatapi.kgrep;

import java.io.InputStream;

public final class KubernetesTestsUtil {

    private KubernetesTestsUtil() {
    }

    public static InputStream getResourceAsStream(String yamlOrJson) {
        return KubernetesTestsUtil.class.getClassLoader().getResourceAsStream(yamlOrJson);
    }
}
