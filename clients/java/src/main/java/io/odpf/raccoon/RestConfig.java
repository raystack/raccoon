package io.odpf.raccoon;

import io.odpf.raccoon.serializer.JsonSerializer;
import io.odpf.raccoon.serializer.Serializer;
import io.odpf.raccoon.wire.JsonWire;
import io.odpf.raccoon.wire.WireMarshaler;
import lombok.Builder;
import lombok.Getter;
import lombok.Setter;
import lombok.Singular;
import org.apache.http.impl.client.CloseableHttpClient;

import java.util.Map;

@Getter
@Setter
@Builder(toBuilder = true)
public class RestConfig {
    private static final int MAX_RETRY = 3;
    private static final int RETRY_WAIT_MILISECONDS = 1000;

    @Builder.Default
    private String url = "http://localhost:8080/api/v1/events";

    @Builder.Default
    private Integer retryMax = MAX_RETRY;

    @Builder.Default
    private long retryWait = RETRY_WAIT_MILISECONDS;

    @Builder.Default
    private Serializer serializer = new JsonSerializer();

    @Builder.Default
    private WireMarshaler marshaler = new JsonWire();

    @Singular
    private Map<String, String> headers;

    private CloseableHttpClient httpClient;
}
