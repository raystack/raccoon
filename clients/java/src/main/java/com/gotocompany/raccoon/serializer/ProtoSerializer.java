package com.gotocompany.raccoon.serializer;

import com.google.protobuf.GeneratedMessageV3;
import com.gotocompany.raccoon.exception.SerializationException;

public class ProtoSerializer implements Serializer {

    /**
     * Converts a proto object to the byte array.
     *
     * @param any proto object to be serialized.
     * @return the proto byte array.
     * @throws SerializationException Throws exceptions if the object is non-proto.
     */
    @Override
    public byte[] serialize(Object any) throws SerializationException {
        if (any instanceof GeneratedMessageV3) {
            return ((GeneratedMessageV3) any).toByteArray();
        }

        throw new SerializationException("Error: unable to serialize non proto object");
    }
}
