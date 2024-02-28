package com.thegreatapi.kgrep;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import jakarta.enterprise.inject.Produces;

class ObjectMapperProducer {

    @Produces
    ObjectMapper produceObjectMapper() {
        return new ObjectMapper(YAMLFactory.builder().build());
    }
}
