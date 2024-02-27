package com.thegreatapi.kgrep;

import jakarta.enterprise.context.ApplicationScoped;

import java.util.Arrays;
import java.util.List;

@ApplicationScoped
class Grep {

    List<String> run(String[] lines, String pattern) {
        return Arrays.stream(lines)
                .filter(line -> line.contains(pattern))
                .toList();
    }
}
