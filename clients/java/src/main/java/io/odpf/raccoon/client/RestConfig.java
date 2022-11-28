package io.odpf.raccoon.client;

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

/*
 * Rest config options for the Raccoon http request.
 */
public class RestConfig {
    /**
     * The default max retry attempts.
     */
    private static final int MAX_RETRY = 3;

    /**
     * The default delay between retries.
     */
    private static final int RETRY_WAIT_MILISECONDS = 1000;

    /**
     * The raccoon http endpoint to make the request.
     */
    @Builder.Default
    private String url = "http://localhost:8080/api/v1/events";

    /**
     * The max retry attempts to be made on failed requests.
     */
    @Builder.Default
    private Integer retryMax = MAX_RETRY;

    /**
     * The sleep time between the retry requests in miliseconds.
     */
    @Builder.Default
    private long retryWait = RETRY_WAIT_MILISECONDS;

    /**
     * The serializer for the raccoon batch request.
     */
    @Builder.Default
    private Serializer serializer = new JsonSerializer();

    /**
     * The marshaler for the wire request made to the raccoon.
     */
    @Builder.Default
    private WireMarshaler marshaler = new JsonWire();

    /**
     * The http request headers.
     */
    @Singular
    private Map<String, String> headers;

    /**
     * The http client.
     */
    private CloseableHttpClient httpClient;
}
