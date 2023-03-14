package com.gotocompany.raccoon.serializer;

import com.google.gson.Gson;
import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.MessageOrBuilder;
import com.google.protobuf.util.JsonFormat;
import com.gotocompany.raccoon.exception.SerializationException;

/**
    Json class to convert raccoon messages.
**/
public class JsonSerializer implements Serializer {

    private final Gson gson;

    /**
     * Initializes the gson instance.
     */
    public JsonSerializer() {
        this.gson = new Gson();
    }

    /**
     * Json serializer for the raccoon messages.
     *
     * @see Serializer#serialize(java.lang.Object)
     * @param any is the object to be serialized.
     * @return returns the json bytes.
     **/
    @Override
    public byte[] serialize(Object any) throws SerializationException {
        if (any instanceof MessageOrBuilder) {
            try {
                return JsonFormat.printer()
                        .omittingInsignificantWhitespace()
                        .preservingProtoFieldNames()
                        .print((MessageOrBuilder) any)
                        .getBytes();
            } catch (InvalidProtocolBufferException e) {
                throw new SerializationException(e.getMessage());
            }
        }

        return this.gson.toJson(any).getBytes();
    }
}
