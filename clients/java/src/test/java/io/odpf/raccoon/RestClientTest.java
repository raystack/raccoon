package io.odpf.raccoon;

import com.github.tomakehurst.wiremock.junit.WireMockRule;
import com.google.protobuf.Timestamp;
import com.google.protobuf.util.JsonFormat;
import io.odpf.proton.raccoon.SendEventResponse;
import io.odpf.raccoon.client.RaccoonClient;
import io.odpf.raccoon.client.RaccoonClientFactory;
import io.odpf.raccoon.client.RestConfig;
import io.odpf.raccoon.model.Event;
import io.odpf.raccoon.model.Response;
import io.odpf.raccoon.model.ResponseCode;
import io.odpf.raccoon.model.ResponseStatus;
import io.odpf.raccoon.serializer.JsonSerializer;
import io.odpf.raccoon.serializer.ProtoSerializer;
import io.odpf.raccoon.wire.JsonWire;
import io.odpf.raccoon.wire.ProtoWire;
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
                                .setStatus(io.odpf.proton.raccoon.Status.STATUS_SUCCESS)
                                .setCode(io.odpf.proton.raccoon.Code.CODE_OK)
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
                                .setStatus(io.odpf.proton.raccoon.Status.STATUS_SUCCESS)
                                .setCode(io.odpf.proton.raccoon.Code.CODE_OK)
                                .setSentTime(TimeUnit.DAYS.toDays(1))
                                .putAllData(data)
                                .build();

                return sendEventResponse.toByteArray();
        }
}
