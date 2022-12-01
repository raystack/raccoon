package io.odpf.raccoon.wire;

import com.google.gson.Gson;
import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.Message;
import com.google.protobuf.MessageOrBuilder;
import com.google.protobuf.util.JsonFormat;
import io.odpf.raccoon.exception.DeserializationException;
import io.odpf.raccoon.exception.SerializationException;

/**
 * Json class to convert raccoon messages to/from the proto3.
 */
public class JsonWire implements WireMarshaler {

    private final Gson gson;

    public JsonWire() {
        this.gson = new Gson();
    }

    /**
     * Converts a raccoon message to the json.
     *
     * @param any object to be marshaled.
     * @return the marshal byte array.
     * @throws SerializationException Throws exceptions if fail to serialize.
     */
    @Override
    public byte[] marshal(Object any) throws SerializationException {
        try {
            if (any instanceof MessageOrBuilder) {
                return JsonFormat.printer()
                        .omittingInsignificantWhitespace()
                        .preservingProtoFieldNames()
                        .print((MessageOrBuilder) any)
                        .getBytes();
            }

            return this.gson.toJson(any).getBytes();

        } catch (InvalidProtocolBufferException e) {
            throw new SerializationException(e.getMessage());
        }
    }

    /**
     * Converts the json bytes back into the object.
     *
     * @param any  object to be unmarshal.
     * @param type object to be used for unmarshal.
     * @return the unmarshal object.
     * @throws DeserializationException Throws exceptions if fail to deserialize.
     */
    @Override
    public Object unmarshal(byte[] any, Object type) throws DeserializationException {
        try {
            if (type instanceof Message.Builder) {
                Message.Builder builder = (Message.Builder) type;
                JsonFormat.parser().merge(new String(any), builder);
                return builder.build();
            }

            return this.gson.fromJson(new String(any), type.getClass());
        } catch (InvalidProtocolBufferException e) {
            throw new DeserializationException(e.getMessage());
        }
    }

    /**
     * Wire content type for the json marshaled request.
     *
     * @return the content type.
     */
    @Override
    public String getContentType() {
        return "application/json";
    }

}
