package io.odpf.raccoon;

import java.util.Map;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@NoArgsConstructor
@AllArgsConstructor
@Data
public class Response {
    private int status;
    private int code;
    private long sentTime;
    private Map<String, String> data;
}
