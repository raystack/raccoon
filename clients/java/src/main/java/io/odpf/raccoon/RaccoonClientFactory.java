package io.odpf.raccoon;

import lombok.NonNull;

/**
 * Factory for raccoon client.
 */
public class RaccoonClientFactory {
    /**
     * Creates the raccoon Rest client.
     *
     * @param restConfig The rest config options.
     * @return The rest client.
     */
    public static RaccoonClient getRestClient(@NonNull RestConfig restConfig) {
        return new RestClient(restConfig);
    }
}
