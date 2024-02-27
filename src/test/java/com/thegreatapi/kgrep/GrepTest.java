package com.thegreatapi.kgrep;

import io.quarkus.test.junit.QuarkusTest;
import jakarta.inject.Inject;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;

@QuarkusTest
class GrepTest {

    @Inject
    Grep grep;

    @Test
    void run() {
        String[] lines = new String[]{
                "one line",
                "two lines",
                "a line ending with kgrep",
                "kgrep in the beginning"
        };

        List<String> occurrences = grep.run(lines, "kgrep");

        assertThat(occurrences).containsExactlyInAnyOrder(
                "a line ending with kgrep",
                "kgrep in the beginning");
    }
}