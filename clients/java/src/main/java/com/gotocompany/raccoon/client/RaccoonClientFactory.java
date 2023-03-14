package com.gotocompany.raccoon.client;

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
