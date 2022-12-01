package io.odpf.raccoon.exception;

/**
 * Exception thrown when the deserialization process fails.
 */
public class DeserializationException extends RuntimeException {
    /**
     * Constructs a new {@link DeserializationException} with specified detail
     * message.
     *
     * @param msg The error message.
     */
    public DeserializationException(final String msg) {
        super(msg);
    }
}
