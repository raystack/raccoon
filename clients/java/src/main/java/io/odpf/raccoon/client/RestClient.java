package io.odpf.raccoon.client;

import com.google.api.client.http.HttpStatusCodes;
import com.google.common.io.ByteStreams;
import com.google.protobuf.ByteString;
import io.odpf.proton.raccoon.Event.Builder;
import io.odpf.proton.raccoon.SendEventRequest;
import io.odpf.proton.raccoon.SendEventResponse;
import io.odpf.raccoon.model.Event;
import io.odpf.raccoon.model.Response;
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
 class RestClient implements RaccoonClient {

    private static final Logger LOGGER = LoggerFactory.getLogger(RestClient.class);

    private final RestConfig restConfig;

    /**
     * @param restConfig The rest config options.
     */
    RestClient(RestConfig restConfig) {
        this.restConfig = restConfig;
        this.restConfig.setHttpClient(this.getHttpClient());
    }

    /**
     * Send creates a batch request, produces a request guid, and sends an HTTP
     * request.
     *
     * @param events The raccoon message array.
     * @return {@link Response} The response for the given raccoon request.
     */
    @Override
    public Response send(Event[] events) {
        String reqGuid = UUID.randomUUID().toString();
        CloseableHttpResponse response = null;
        String errorMessage;
        try {
            LOGGER.info("send: started request, url:{}, req-id: {}", this.restConfig.getUrl(), reqGuid);
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

            response = this.restConfig.getHttpClient().execute(httPost);
            SendEventResponse eventResponse = (SendEventResponse) this.restConfig.getMarshaler().unmarshal(
                    ByteStreams.toByteArray(response.getEntity()
                            .getContent()),
                    this.restConfig.getMarshaler() instanceof JsonWire ? SendEventResponse.newBuilder()
                            : SendEventResponse.parser());

            return Response.builder().isSuccess(true)
                    .reqGuid(reqGuid)
                    .status(eventResponse.getStatusValue())
                    .code(eventResponse.getCodeValue())
                    .sentTime(eventResponse.getSentTime())
                    .data(eventResponse.getDataMap())
                    .build();

        } catch (Exception ex) {
            LOGGER.error(ex.getMessage());
            errorMessage = ex.getMessage();
        } finally {
            if (response != null) {
                try {
                    response.close();
                } catch (IOException e) {
                    LOGGER.error("send: exception when closing httpclient", e);
                    errorMessage = e.getMessage();
                }
            }
            LOGGER.info("send: ended request, url:{}, req-id: {}", this.restConfig.getUrl(), reqGuid);
        }

        return Response.builder()
                .isSuccess(false)
                .reqGuid(reqGuid)
                .errorMessage(errorMessage)
                .build();
    }

    /**
     * Creates the new HTTP client and set the retry handler with the provided settings.
     *
     * @return The new HTTP client.
     */
    private CloseableHttpClient getHttpClient() {

        return HttpClients
                .custom()
                .setRetryHandler((exception, executionCount, context) -> executionCount < restConfig.getRetryMax())
                .setServiceUnavailableRetryStrategy(new ServiceUnavailableRetryStrategy() {

                    @Override
                    public boolean retryRequest(HttpResponse response, int executionCount, HttpContext context) {
                        if (executionCount < restConfig.getRetryMax() && response.getStatusLine().getStatusCode() != HttpStatusCodes.STATUS_CODE_OK) {
                            LOGGER.warn("Retrying the http request, retry-count:{}, response-code:{}", executionCount, response.getStatusLine().getStatusCode());
                            return true;
                        }

                        return false;
                    }

                    @Override
                    public long getRetryInterval() {
                        return restConfig.getRetryWait();
                    }

                })
                .build();
    }
}
