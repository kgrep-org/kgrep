package com.thegreatapi.kgrep.grep;

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
                "one text",
                "two lines",
                "a text ending with kgrep",
                "kgrep in the beginning"
        };

        List<Occurrence> occurrences = grep.run(lines, "kgrep");

        assertThat(occurrences).containsExactlyInAnyOrder(
                new Occurrence(3, "a text ending with kgrep"),
                new Occurrence(4, "kgrep in the beginning"));
    }
}