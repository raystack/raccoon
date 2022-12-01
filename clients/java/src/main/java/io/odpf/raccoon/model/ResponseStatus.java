package io.odpf.raccoon.model;

public class ResponseStatus {
    /*
     * `UNSPECIFIED` indicates if request is failed for the unknown reasons.
     */
    public static final int STATUS_UNSPECIFIED = 0;

    /*
     * 'SUCCESS' indicates if the request is successful.
     */
    public static final int STATUS_SUCCESS = 1;

    /*
     * 'ERROR' indicates if the request is failed.
     */
    public static final int STATUS_ERROR = 2;
}
