package io.odpf.raccoon;

import com.google.api.client.http.HttpStatusCodes;
import com.google.protobuf.ByteString;
import io.odpf.proton.raccoon.Event.Builder;
import io.odpf.proton.raccoon.SendEventRequest;
import io.odpf.proton.raccoon.SendEventResponse;
import io.odpf.raccoon.wire.JsonWire;
import org.apache.http.HttpResponse;
import org.apache.http.client.ServiceUnavailableRetryStrategy;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.ByteArrayEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.protocol.HttpContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.UUID;

/**
 * Class for the http client implementation.
 */
public class RestClient implements RaccoonClient {

    private static final Logger LOGGER = LoggerFactory.getLogger(RestClient.class);

    private final RestConfig restConfig;

    /**
     * @param restConfig The rest config options.
     */
    public RestClient(RestConfig restConfig) {
        this.restConfig = restConfig;
        this.restConfig.setHttpClient(this.getHttpClient());
    }

    /**
     * Send makes the http request to the raccoon service.
     *
     * @param events The raccoon event message array.
     * @return The response for the given raccoon event request.
     * @throws Exception Throws the exception if the request fails.
     */
    @Override
    public Response send(Event[] events) throws Exception {
        String reqGuid = UUID.randomUUID().toString();
        LOGGER.info("started request, url:{}, req-id: {}", this.restConfig.getUrl(), reqGuid);
        SendEventRequest.Builder builder = SendEventRequest.newBuilder();
        builder.setReqGuid(reqGuid);

        for (Event e : events) {
            Builder protoEvents = io.odpf.proton.raccoon.Event.newBuilder();
            protoEvents.setType(e.getType());
            protoEvents.setEventBytes(ByteString.copyFrom(this.restConfig.getSerializer().serialize(e.getData())));

            builder.addEvents(protoEvents.build());
        }

        // wire request.
        SendEventRequest postRequest = builder.build();

        HttpPost httPost = new HttpPost(restConfig.getUrl());
        httPost.addHeader("Content-Type", this.restConfig.getMarshaler().getContentType());
        httPost.setEntity(new ByteArrayEntity(this.restConfig.getMarshaler().marshal(postRequest)));

        this.restConfig.getHeaders().forEach(httPost::addHeader);

        CloseableHttpResponse response;

        response = this.restConfig.getHttpClient().execute(httPost);
        try {
            if (response.getStatusLine().getStatusCode() == HttpStatusCodes.STATUS_CODE_OK) {
                SendEventResponse eventResponse = (SendEventResponse) this.restConfig.getMarshaler().unmarshal(
                        response.getEntity()
                                .getContent()
                                .readAllBytes(),
                        this.restConfig.getMarshaler() instanceof JsonWire ? SendEventResponse.newBuilder()
                                : SendEventResponse.parser());

                return new Response(
                        eventResponse.getStatusValue(),
                        eventResponse.getCodeValue(),
                        eventResponse.getSentTime(),
                        eventResponse.getDataMap());
            }
        } catch (IOException e1) {
            e1.printStackTrace();
            LOGGER.error(e1.getMessage());
        } finally {
            response.close();
            LOGGER.info("started ended, url:{}, req-id: {}", this.restConfig.getUrl(), reqGuid);
        }

        return null;
    }

    private CloseableHttpClient getHttpClient() {

        return HttpClients
                .custom()
                .setRetryHandler((exception, executionCount, context) -> executionCount < restConfig.getRetryMax())
                .setServiceUnavailableRetryStrategy(new ServiceUnavailableRetryStrategy() {

                    @Override
                    public boolean retryRequest(HttpResponse response, int executionCount, HttpContext context) {
                        LOGGER.warn("Retrying the http request, retry-count:{}, response-code:{}", executionCount,
                                response.getStatusLine().getStatusCode());
                        return executionCount < restConfig.getRetryMax()
                                && response.getStatusLine().getStatusCode() != HttpStatusCodes.STATUS_CODE_OK;
                    }

                    @Override
                    public long getRetryInterval() {
                        return restConfig.getRetryWait();
                    }

                })
                .build();
    }
}
