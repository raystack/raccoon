package com.gotocompany.raccoon.client;

import com.gotocompany.raccoon.model.Event;
import com.gotocompany.raccoon.model.Response;

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
