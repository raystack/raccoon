package io.odpf.raccoon.client;

import io.odpf.raccoon.model.Event;
import io.odpf.raccoon.model.Response;

/**
 * An interface for the raccoon clients.
 */
public interface RaccoonClient {
    /**
     * Sends a request to raccoon with the message provided.
     *
     * @param events The raccoon event message array.
     * @return {@link Response} The raccoon response.
     */
    Response send(Event[] events);
}
