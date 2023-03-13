package com.gotocompany.raccoon.wire;

import com.gotocompany.raccoon.exception.DeserializationException;
import com.gotocompany.raccoon.exception.SerializationException;

/**
 * An interface for marshaling/unmarshalling the wire requests.
 */
public interface WireMarshaler {
    /**
     * Marshal the given object into the byte array.
     *
     * @param any serializable object.
     * @return the byte array.
     * @throws SerializationException throw when the serialization fails.
     */
    byte[] marshal(Object any) throws SerializationException;

    /**
     * Unmarshal the given byte array into the object.
     *
     * @param any serializable byte array.
     * @param type the type to be
     * @return The object un-marshalled from the byte array.
     * @throws DeserializationException throw when the deserialization fails.
     */
    Object unmarshal(byte[] any, Object type) throws DeserializationException;

    /**
     * Content type of the wire request.
     *
     * @return the content type.
     */
    String getContentType();
}
