package io.odpf.raccoon.client;

import io.odpf.raccoon.model.Event;
import io.odpf.raccoon.model.Result;

/**
 * An interface for the raccoon clients.
 */
public interface RaccoonClient {
    /**
     * Sends a request to raccoon with the message provided.
     *
     * @param <T> The response type of the raccoon request.
     *
     * @param events The raccoon event message array.
     * @return The raccoon response.
     */
    <T> Result<T> send(Event[] events);
}
