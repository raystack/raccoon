package io.odpf.raccoon;

/**
 * An interface for the raccoon clients.
 */
public interface RaccoonClient {
    /**
     * Sends a request to raccoon with the message provided.
     *
     * @param events The raccoon event message array.
     * @return The raccoon response.
     * @throws Exception Throws the exception if the request fails.
     */
    Response send(Event[] events) throws Exception;
}
