package com.thegreatapi.kgrep;

import java.io.InputStream;

final class KubernetesTestsUtil {

    private KubernetesTestsUtil() {
    }

    static InputStream getResourceAsStream(String yamlOrJson) {
        return ConfigMapGrepperTest.class.getClassLoader().getResourceAsStream(yamlOrJson);
    }
}
