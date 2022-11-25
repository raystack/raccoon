package io.odpf.raccoon.model;

import lombok.Builder;
import lombok.Data;

@Data
@Builder(toBuilder = true)

/*
 * Result is the wrapper class for the raccoon client response,
 * It holds the status of the request, error message and the response.
 */
public class Result<T> {
    /**
     * The status of the raccoon request.
     */
    private boolean isSuccess;

    /**
     * The request guid of the request.
     */
    private String reqGuid;

    /**
     * The raccoon response.
     */
    private T response;

    /**
     * The error message which gets populated on the failed requests.
     */
    private String errorMessage;
}
