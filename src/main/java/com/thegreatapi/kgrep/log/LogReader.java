package com.thegreatapi.kgrep.log;

interface LogReader {

    String read(String namespace, String podName, String containerName);
}
