package io.odpf.raccoon.serializer;

import com.google.protobuf.Timestamp;
import com.google.protobuf.util.JsonFormat;

import sample.PageEventProto;

import org.junit.Assert;
import org.junit.Test;

import java.util.UUID;

public class JsonSerializerTest {

    @Test
    public void testSerialize() throws Exception {
        // Arrange
        Serializer serializer = new JsonSerializer();

        // Act
        PageEventProto.PageEvent pageEvent = PageEventProto.PageEvent.newBuilder()
                .setEventGuid(UUID.randomUUID().toString())
                .setEventName("clicked")
                .setSentTime(Timestamp.newBuilder().build())
                .build();

        byte[] actualBytes = serializer.serialize(pageEvent);

        PageEventProto.PageEvent.Builder builder = PageEventProto.PageEvent.newBuilder();
        JsonFormat.parser().merge(new String(actualBytes), builder);
        PageEventProto.PageEvent actual = builder.build();

        // Assert
        Assert.assertEquals(pageEvent.getEventGuid(), actual.getEventGuid());
        Assert.assertEquals(pageEvent.getEventName(), actual.getEventName());
    }
}
