package com.gotocompany.raccoon.model;

public class ResponseCode {
    /*
     * `CODE_UNSPECIFIED` indicates no appropriate/existing code can describe it.
     */
    public static final int CODE_UNSPECIFIED = 0;

    /*
     * `OK` indicates the request is processed successfully.
     */
    public static final int CODE_OK = 1;

    /*
     * `BAD_REQUEST` indicates there is something wrong with the request.
     */
    public static final int CODE_BAD_REQUEST = 2;

    /*
     * `INTERNAL_ERROR` indicates that Raccoon encountered an unexpected condition
     * that prevented it from fulfilling the request.
     */
    public static final int CODE_INTERNAL_ERROR = 3;

    /*
     * `MAX_CONNECTION_LIMIT_REACHED` indicates that Raccoon is unable to accepts
     * new connection due to max connection is reached.
     */

    public static final int CODE_MAX_CONNECTION_LIMIT_REACHED = 4;

    /*
     * `MAX_USER_LIMIT_REACHED` indicates that existing connection with the same ID.
     */
    public static final int CODE_MAX_USER_LIMIT_REACHED = 5;
}
