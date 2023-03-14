package com.gotocompany.raccoon.serializer;

import com.gotocompany.raccoon.exception.SerializationException;

/*
    Interface serializer defines a conversion for raccoon message to byte sequence.
*/
public interface Serializer {
    /**
     * Serialize the given object into the byte array.
     *
     * @param any object to be serialized.
     * @return the serialized byte array.
     * @throws SerializationException if the object fail to serialize.
     */
    byte[] serialize(Object any) throws SerializationException;
}
