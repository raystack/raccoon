package io.odpf.raccoon.model;

import lombok.AllArgsConstructor;
import lombok.Data;

@AllArgsConstructor
@Data

/*
 * Class holds the raccoon event and it's type name.
 */
public class Event {

    /**
     * The event type name.
     */
    private String type;

    /**
     * The event batch.
     */
    private Object data;
}
