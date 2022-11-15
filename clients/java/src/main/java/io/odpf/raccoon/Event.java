package io.odpf.raccoon;

import lombok.AllArgsConstructor;
import lombok.Data;

@AllArgsConstructor
@Data
public class Event {
    private String type;
    private Object data;
}
