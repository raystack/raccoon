package com.gotocompany.raccoon.model;

import java.util.Map;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@NoArgsConstructor
@AllArgsConstructor
@Data
@Builder(toBuilder = true)

/*
 * Response holds the information about the request.
 */
public class Response {

  /**
   * The status of the raccoon request.
   */
  private boolean isSuccess;

  /**
   * The request guid of the request.
   */
  private String reqGuid;

  /**
   * The error message which gets populated on the failed requests.
   */
  private String errorMessage;

  /**
   * The {@link ResponseStatus} denotes status of the request.
   */
  private int status;

  /**
   * The {@link ResponseCode} gives more detail of what happened to the request.
   */
  private int code;

  /**
   * The sentTime is UNIX timestamp populated by Raccoon by the time the response
   * is
   * sent.
   */
  private long sentTime;

  /**
   * Data is arbitrary extra metadata.
   * The data map will contain the "req_guid" key
   */
  private Map<String, String> data;
}
