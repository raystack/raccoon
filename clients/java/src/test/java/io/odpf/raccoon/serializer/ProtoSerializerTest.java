package io.odpf.raccoon.serializer;

import com.google.protobuf.Timestamp;

import sample.PageEventProto;

import org.junit.Assert;
import org.junit.Test;

import java.util.UUID;

public class ProtoSerializerTest {

    @Test
    public void testSerialize() throws Exception {
        // Arrange
        Serializer serializer = new ProtoSerializer();

        // Act
        PageEventProto.PageEvent pageEvent = PageEventProto.PageEvent.newBuilder()
                .setEventGuid(UUID.randomUUID().toString())
                .setEventName("clicked")
                .setSentTime(Timestamp.newBuilder().build())
                .build();

        byte[] pageBytes = serializer.serialize(pageEvent);

        PageEventProto.PageEvent actual = PageEventProto.PageEvent.parser().parseFrom(pageBytes);

        // Assert
        Assert.assertEquals(pageEvent.getEventGuid(), actual.getEventGuid());
        Assert.assertEquals(pageEvent.getEventName(), actual.getEventName());
    }
}
