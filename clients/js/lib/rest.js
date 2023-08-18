import axios from 'axios';
import { v4 as uuidv4 } from 'uuid';

import { createJsonSerializer } from './serializer/json_serializer.js';
import { createProtobufSerializer } from './serializer/proto_serializer.js';
import { retry } from './retry/retry.js';
import { createProtoMarshaller } from './wire/proto_wire.js';
import { createJsonMarshaller } from './wire/json_wire.js';
import protos from '../protos/proton_compiled.js';

const NANOSECONDS_PER_MILLISECOND = 1e6;

const SerializationType = Object.freeze({
    JSON: 'json',
    PROTOBUF: 'protobuf'
});

const WireType = Object.freeze({
    JSON: 'application/json',
    PROTOBUF: 'application/proto'
});

class RaccoonClient {
    /**
     * Creates a new instance of the RaccoonClient.
     *
     * @constructor
     * @param {Object} options - Configuration options for the RaccoonClient.
     * @param {string} [options.serializationType] - The serialization type to use, either 'protobuf' or 'json'.
     * @param {Object} [options.wireType] - The wire configuration, containing ContentType.
     * @param {Object} [options.headers] - Custom headers to be included in the HTTP requests.
     * @param {number} [options.retryMax=3] - The maximum number of retry attempts for failed requests.
     * @param {number} [options.retryWait=1000] - The time in milliseconds to wait between retry attempts.
     * @param {string} [options.url=''] - The base URL for the API requests.
     * @param {string} [options.logger=''] - Logger object for logging.
     * @param {number} [options.timeout=5000] - The timeout in milliseconds.
     * @returns {RaccoonClient} A new instance of the RaccoonClient.
     */
    constructor(options = {}) {
        if (!Object.values(SerializationType).includes(options.serializationType)) {
            throw new Error(`Invalid serializationType: ${options.serializationType}`);
        }
        this.serialize = options.serializationType === SerializationType.PROTOBUF
            ? createProtobufSerializer()
            : createJsonSerializer();
        if (!Object.values(WireType).includes(options.wireType)) {
            throw new Error(`Invalid wireType: ${options.wireType}`);
        }
        this.marshaller = options.wireType === WireType.PROTOBUF
            ? createProtoMarshaller()
            : createJsonMarshaller();
        this.headers = options.headers || {};
        this.retryMax = options.retryMax || 3;
        this.retryWait = options.retryWait || 5000;
        this.url = options.url || '';
        this.logger = options.logger || console
        this.timeout = options.timeout || 5000;
        this.uuidGenerator = (() => uuidv4());
        this.httpClient = axios.create();
    }

    /**
     * Sends a batch of events to the server.
     *
     * @async
     * @param {Array} events - An array of event objects to send.
     * @returns {Promise} A promise that resolves to an object containing the request ID, response, and error details (if any).
     * @throws {Error} Throws an error if the event is invalid or if there's an issue with the request.
     */
    async send(events) {
        const requestId = this.uuidGenerator();
        try {
            if (!events || events.length === 0) {
                throw new Error('No events provided');
            }
            this.logger.info(`started request, url: ${this.url}, req-id: ${requestId}`);
            const eventsToSend = [];
            for (const event of events) {
                if (event && event.type && event.data) {
                    let eventMessage = new protos.raystack.raccoon.v1beta1.Event();
                    eventMessage.type = event.type;
                    eventMessage.event_bytes = this.serialize(event.data);
                    eventsToSend.push(eventMessage);
                } else {
                    throw new Error(`Invalid event: ${JSON.stringify(event)}`);
                }
            }
            const sendEventRequest = new protos.raystack.raccoon.v1beta1.SendEventRequest();
            sendEventRequest.req_guid = requestId;

            const now = Date.now();
            sendEventRequest.sent_time = protos.google.protobuf.Timestamp.create({
                seconds: Math.floor(now / 1000),
                nanos: (now % 1000) * NANOSECONDS_PER_MILLISECOND
            });
            sendEventRequest.events = eventsToSend;

            const response = await retry(
                async () => this.executeRequest(this.marshaller.marshal(sendEventRequest)),
                this.retryMax,
                this.retryWait
            );

            const sendEventResponse = this.marshaller.unmarshal(response, protos.raystack.raccoon.v1beta1.SendEventResponse);

            this.logger.info(`ended request, url: ${this.url}, req-id: ${requestId}`);
            return {
                reqId: requestId,
                response: sendEventResponse.toJSON(),
                error: null,
            };
        } catch (error) {
            this.logger.error(`error, url: ${this.url}, req-id: ${requestId}, ${error}`);
            throw new Error(`req-id: ${requestId}, error: ${error}`);
        }
    }

    async executeRequest(raccoonRequest) {
        this.headers['Content-Type'] = this.marshaller.getContentType();
        const response = await this.httpClient.post(this.url, raccoonRequest, {
            headers: this.headers,
            timeout: this.timeout,
            responseType: this.marshaller.getResponseType()
        });
        return response.data;
    }
}

export { RaccoonClient, SerializationType, WireType };
