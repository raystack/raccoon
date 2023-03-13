package com.gotocompany.raccoon.wire;

import com.google.protobuf.GeneratedMessageV3;
import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.Parser;
import com.gotocompany.raccoon.exception.DeserializationException;
import com.gotocompany.raccoon.exception.SerializationException;

/**
 * Proto class to convert raccoon messages to/from the proto3.
 */
public class ProtoWire implements WireMarshaler {

    /**
     * Converts a raccoon request to the proto.
     *
     * @param any proto object to be marshaled.
     * @return the marshal byte array.
     * @throws SerializationException Throws when the serialization fails.
     */
    @Override
    public byte[] marshal(Object any) throws SerializationException {
        if (any instanceof GeneratedMessageV3) {
            return ((GeneratedMessageV3) any).toByteArray();
        }

        throw new SerializationException("unable to serialize the non proto message");
    }

    /**
     * Converts the raccoon request proto bytes back into the object.
     *
     * @param any  proto object to be unmarshal.
     * @param type prototype to be used for unmarshal.
     * @return the unmarshalled object.
     * @throws DeserializationException Throws when the deserialization fails.
     */
    @Override
    public Object unmarshal(byte[] any, Object type) throws DeserializationException {
        try {
            if (type instanceof Parser<?>) {
                return ((Parser<?>) type).parseFrom(any);
            }
        } catch (InvalidProtocolBufferException e) {
            throw new DeserializationException(e.getMessage());
        }

        throw new DeserializationException("unable to deserialize using the non proto parser");
    }

    /**
     * Wire content type for the proto marshaled request.
     *
     * @return the content type.
     */
    @Override
    public String getContentType() {
        return "application/proto";
    }

}
