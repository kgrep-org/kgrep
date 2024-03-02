package com.thegreatapi.kgrep.grep;

import jakarta.enterprise.context.ApplicationScoped;

import java.util.ArrayList;
import java.util.List;

@ApplicationScoped
public final class Grep {

    public List<Occurrence> run(String[] lines, String pattern) {
        List<Occurrence> list = new ArrayList<>();
        for (int i = 0; i < lines.length; i++) {
            String line = lines[i];
            if (line.contains(pattern)) {
                list.add(new Occurrence(i + 1, line));
            }
        }
        return list;
    }
}
