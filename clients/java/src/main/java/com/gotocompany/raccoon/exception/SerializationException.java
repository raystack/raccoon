package com.gotocompany.raccoon.exception;

/**
 * Exception thrown when the serialization process fails.
 */
public class SerializationException extends RuntimeException {
    /**
     * Constructs a new {@link SerializationException} with specified detail
     * message.
     *
     * @param msg The error message.
     */
    public SerializationException(final String msg) {
        super(msg);
    }
}
