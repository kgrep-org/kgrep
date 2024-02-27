package com.thegreatapi.kgrep;

interface LogReader {

    String read(String namespace, String podName, String containerName);
}
