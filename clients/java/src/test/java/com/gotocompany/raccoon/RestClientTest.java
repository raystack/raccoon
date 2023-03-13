package com.gotocompany.raccoon;

import com.github.tomakehurst.wiremock.junit.WireMockRule;
import com.google.protobuf.Timestamp;
import com.google.protobuf.util.JsonFormat;
import com.gotocompany.proton.raccoon.Code;
import com.gotocompany.proton.raccoon.Status;
import com.gotocompany.raccoon.client.RaccoonClient;
import com.gotocompany.raccoon.client.RaccoonClientFactory;
import com.gotocompany.raccoon.client.RestConfig;
import com.gotocompany.raccoon.model.Event;
import com.gotocompany.raccoon.model.Response;
import com.gotocompany.raccoon.model.ResponseCode;
import com.gotocompany.raccoon.model.ResponseStatus;
import com.gotocompany.raccoon.serializer.JsonSerializer;
import com.gotocompany.raccoon.serializer.ProtoSerializer;
import com.gotocompany.proton.raccoon.SendEventResponse;
import com.gotocompany.raccoon.wire.JsonWire;
import com.gotocompany.raccoon.wire.ProtoWire;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import sample.PageEventProto;

import java.util.HashMap;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

import static com.github.tomakehurst.wiremock.client.WireMock.*;

public class RestClientTest {

        @Rule
        public WireMockRule service = new WireMockRule(8082);

        @Test
        public void testProtoSend() throws Exception {

                String reqGuid = UUID.randomUUID().toString();
                service.stubFor(
                                post("/api/v1/events")
                                                .withHeader("Content-Type", containing("proto"))
                                                .willReturn(ok()
                                                                .withHeader("Content-Type", "application/proto")
                                                                .withBody(getProtoResponse(reqGuid))));

                RestConfig config = RestConfig.builder()
                                .url(service.url("/api/v1/events"))
                                .header("x-connection-id", "123")
                                .serializer(new ProtoSerializer())
                                .marshaler(new ProtoWire()).build();

                RaccoonClient restClient = RaccoonClientFactory.getRestClient(config);

                PageEventProto.PageEvent pageEvent = getPageEvent(reqGuid);

                Response response = restClient.send(
                                new Event[] {
                                                new Event("page", pageEvent)
                                });

                Assert.assertTrue(response.isSuccess());
                Assert.assertEquals(response.getCode(), ResponseCode.CODE_OK);
                Assert.assertEquals(response.getStatus(), ResponseStatus.STATUS_SUCCESS);
                Assert.assertTrue(response.getData().containsKey("req_guid"));
                Assert.assertNull(response.getErrorMessage());
                Assert.assertNotNull(response.getReqGuid());
        }

        @Test
        public void testJsonSend() throws Exception {
                String reqGuid = UUID.randomUUID().toString();
                service.stubFor(
                                post("/api/v1/events")
                                                .withHeader("Content-Type", containing("json"))
                                                .willReturn(ok()
                                                                .withHeader("Content-Type", "application/json")
                                                                .withBody(getJsonResponse(reqGuid))));

                RestConfig config = RestConfig.builder()
                                .url(service.url("/api/v1/events"))
                                .header("x-connection-id", "123")
                                .serializer(new JsonSerializer())
                                .marshaler(new JsonWire()).build();

                RaccoonClient restClient = RaccoonClientFactory.getRestClient(config);

                PageEventProto.PageEvent pageEvent = getPageEvent(reqGuid);

                Response response = restClient.send(
                                new Event[] {
                                                new Event("page", pageEvent)
                                });

                Assert.assertTrue(response.isSuccess());
                Assert.assertEquals(response.getCode(), ResponseCode.CODE_OK);
                Assert.assertEquals(response.getStatus(), ResponseStatus.STATUS_SUCCESS);
                Assert.assertTrue(response.getData().containsKey("req_guid"));
                Assert.assertNull(response.getErrorMessage());
                Assert.assertNotNull(response.getReqGuid());
        }

        @Test
        public void testSendServiceUnavailable() {
                String reqGuid = UUID.randomUUID().toString();
                service.stubFor(
                                post("/api/v1/events")
                                                .withHeader("Content-Type", containing("json"))
                                                .willReturn(serviceUnavailable()
                                                                .withHeader("Content-Type", "application/json")));

                RestConfig config = RestConfig.builder()
                                .url(service.url("/api/v1/events"))
                                .header("x-connection-id", "123")
                                .serializer(new JsonSerializer())
                                .marshaler(new JsonWire()).build();

                RaccoonClient restClient = RaccoonClientFactory.getRestClient(config);

                PageEventProto.PageEvent pageEvent = getPageEvent(reqGuid);

                Response response = restClient.send(
                                new Event[] {
                                                new Event("page", pageEvent)
                                });

                Assert.assertFalse(response.isSuccess());
                Assert.assertNull(response.getData());
                Assert.assertNotNull(response.getErrorMessage());
                Assert.assertNotNull(response.getReqGuid());
        }

        private static PageEventProto.PageEvent getPageEvent(String reqGuid) {
                return PageEventProto.PageEvent.newBuilder()
                                .setEventGuid(reqGuid)
                                .setEventName("clicked")
                                .setSentTime(Timestamp.newBuilder().build())
                                .build();
        }

        private static String getJsonResponse(String reqGuid) throws Exception {
                HashMap<String, String> data = new HashMap<>();
                data.put("req_guid", reqGuid);

                SendEventResponse sendEventResponse = SendEventResponse.newBuilder()
                                .setStatus(Status.STATUS_SUCCESS)
                                .setCode(Code.CODE_OK)
                                .setSentTime(TimeUnit.DAYS.toDays(1))
                                .putAllData(data)
                                .build();

                return JsonFormat.printer()
                                .omittingInsignificantWhitespace()
                                .preservingProtoFieldNames()
                                .print(sendEventResponse);
        }

        private static byte[] getProtoResponse(String reqGuid) {
                HashMap<String, String> data = new HashMap<>();
                data.put("req_guid", reqGuid);

                SendEventResponse sendEventResponse = SendEventResponse.newBuilder()
                                .setStatus(Status.STATUS_SUCCESS)
                                .setCode(Code.CODE_OK)
                                .setSentTime(TimeUnit.DAYS.toDays(1))
                                .putAllData(data)
                                .build();

                return sendEventResponse.toByteArray();
        }
}
