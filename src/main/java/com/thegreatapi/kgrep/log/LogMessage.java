package com.thegreatapi.kgrep.log;

record LogMessage(String podName, String containerName, String message, int lineNumber) {
}
