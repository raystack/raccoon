package io.odpf.raccoon;

import io.odpf.raccoon.client.RaccoonClient;
import io.odpf.raccoon.client.rest.RestClient;
import io.odpf.raccoon.client.rest.RestConfig;
import lombok.NonNull;

/**
 * Factory for raccoon client.
 */
public class RaccoonClientFactory {
    /**
     * Creates the new Rest client for the raccoon.
     *
     * @param restConfig The rest config options.
     * @return RestClient The rest client.
     */
    public static RaccoonClient getRestClient(@NonNull RestConfig restConfig) {
        return new RestClient(restConfig);
    }
}
