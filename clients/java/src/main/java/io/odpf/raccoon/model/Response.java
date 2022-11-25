package io.odpf.raccoon.model;

import java.util.Map;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@NoArgsConstructor
@AllArgsConstructor
@Data

/*
 * Response holds the information about the request.
 */
public class Response {

  /**
   * status denotes status of the request.
   */
  private int status;

  /**
   * code gives more detail of what happened to the request.
   */
  private int code;

  /**
   * sentTime is UNIX timestamp populated by Raccoon by the time the response is
   * sent.
   */
  private long sentTime;

  /**
   * data is arbitrary extra metadata.
   */
  private Map<String, String> data;
}
