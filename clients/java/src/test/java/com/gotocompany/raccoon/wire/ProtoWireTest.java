package com.gotocompany.raccoon.wire;

import com.google.protobuf.Timestamp;
import com.gotocompany.proton.raccoon.Event;
import com.gotocompany.proton.raccoon.SendEventRequest;
import sample.PageEventProto;

import org.junit.Assert;
import org.junit.Test;

import java.util.UUID;

public class ProtoWireTest {
        private final WireMarshaler marshaler;

        public ProtoWireTest() {
                this.marshaler = new ProtoWire();
        }

        @Test
        public void testMarshal() throws Exception {
                // arrange
                PageEventProto.PageEvent pageEvent = PageEventProto.PageEvent.newBuilder()
                                .setEventGuid(UUID.randomUUID().toString())
                                .setEventName("clicked")
                                .setSentTime(Timestamp.newBuilder().build())
                                .build();

                SendEventRequest sendEventRequest = SendEventRequest.newBuilder()
                                .setReqGuid(UUID.randomUUID().toString())
                                .addEvents(
                                                Event.newBuilder()
                                                                .setType("page-event")
                                                                .setEventBytes(pageEvent.toByteString())
                                                                .build())
                                .build();

                // act
                byte[] requestBytes = this.marshaler.marshal(sendEventRequest);
                SendEventRequest actual = SendEventRequest.parseFrom(requestBytes);

                // assert
                Assert.assertEquals(sendEventRequest.getReqGuid(), actual.getReqGuid());
                Assert.assertEquals(sendEventRequest.getEventsCount(), actual.getEventsCount());
                Assert.assertEquals(sendEventRequest.getEvents(0).getEventBytes(), actual.getEvents(0).getEventBytes());
        }

        @Test
        public void testUnmarshal() throws Exception {
                // arrange
                PageEventProto.PageEvent pageEvent = PageEventProto.PageEvent.newBuilder()
                                .setEventGuid(UUID.randomUUID().toString())
                                .setEventName("clicked")
                                .setSentTime(Timestamp.newBuilder().build())
                                .build();

                SendEventRequest sendEventRequest = SendEventRequest.newBuilder()
                                .setReqGuid(UUID.randomUUID().toString())
                                .addEvents(
                                                Event.newBuilder()
                                                                .setType("page-event")
                                                                .setEventBytes(pageEvent.toByteString())
                                                                .build())
                                .build();

                // act
                byte[] requestBytes = this.marshaler.marshal(sendEventRequest);
                SendEventRequest actual = (SendEventRequest) this.marshaler.unmarshal(requestBytes,
                                SendEventRequest.parser());

                // assert
                Assert.assertEquals(sendEventRequest.getReqGuid(), actual.getReqGuid());
                Assert.assertEquals(sendEventRequest.getEventsCount(), actual.getEventsCount());
                Assert.assertEquals(sendEventRequest.getEvents(0).getEventBytes(), actual.getEvents(0).getEventBytes());
        }

        @Test
        public void testGetContentType() {
                Assert.assertEquals("application/proto", this.marshaler.getContentType());
        }
}
