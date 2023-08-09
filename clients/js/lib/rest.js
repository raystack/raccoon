import axios from 'axios';
import { v4 as uuidv4 } from 'uuid';

import { createJsonSerializer } from './serializer/json_serializer.js';
import { createProtobufSerializer } from './serializer/proto_serializer.js';
import { retry } from './retry/retry.js';

const NANOSECONDS_PER_MILLISECOND = 1e6;

class RaccoonClient {
    /**
     * Creates a new instance of the RaccoonClient.
     *
     * @constructor
     * @param {Object} options - Configuration options for the RaccoonClient.
     * @param {string} [options.serializationType='json'] - The serialization type to use, either 'protobuf' or 'json'.
     * @param {Object} [options.wire] - The wire configuration, containing options like ContentType.
     * @param {Object} [options.headers] - Custom headers to be included in the HTTP requests.
     * @param {number} [options.retryMax=3] - The maximum number of retry attempts for failed requests.
     * @param {number} [options.retryWait=1000] - The time in milliseconds to wait between retry attempts.
     * @param {string} [options.url=''] - The base URL for the API requests.
     * @param {string} [options.logger=''] - Logger object for logging.
     * @returns {RaccoonClient} A new instance of the RaccoonClient.
     */
    constructor(options = {}, httpClient, uuidGenerator) {
        this.serialize = options.serializationType === 'protobuf'
            ? createProtobufSerializer()
            : createJsonSerializer();
        this.wire = options.wire || { ContentType: 'application/json' };
        this.httpClient = httpClient || axios.create();
        this.headers = options.headers || {};
        this.retryMax = options.retryMax || 3;
        this.retryWait = options.retryWait || 5000;
        this.url = options.url || '';
        this.logger = options.logger || console
        this.uuidGenerator = uuidGenerator || (() => uuidv4());
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
                    eventsToSend.push({
                        type: event.type,
                        event_bytes: this.serialize(event.data),
                    });
                } else {
                    throw new Error(`Invalid event: ${JSON.stringify(event)}`);
                }
            }

            const sentTime = getCurrentTimestamp();

            const raccoonRequest = {
                req_guid: requestId,
                events: eventsToSend,
                sent_time: sentTime,
            };

            const response = await retry(
                async () => this.executeRequest(JSON.stringify(raccoonRequest)),
                this.retryMax,
                this.retryWait
            );

            const status = parseInt(response.status, 10);
            const code = parseInt(response.code, 10);

            this.logger.info(`ended request, url: ${this.url}, req-id: ${requestId}`);
            return {
                reqId: requestId,
                response: {
                    Status: status,
                    Code: code,
                    SentTime: response.sent_time,
                    Data: response.data,
                },
                error: null,
            };
        } catch (error) {
            this.logger.error(`error, url: ${this.url}, req-id: ${requestId}, ${error}`);
            throw new Error(`req-id: ${requestId}, error: ${error}`);
        }
    }

    async executeRequest(raccoonRequest) {
        this.headers['Content-Type'] = this.wire.ContentType;
        const response = await this.httpClient.post(this.url, raccoonRequest, {
            headers: this.headers,
        });
        return response.data;
    }
}

function getCurrentTimestamp() {
    const now = Date.now();
    return ({
        seconds: Math.floor(now / 1000),
        nanos: (now % 1000) * NANOSECONDS_PER_MILLISECOND
    });
}

export { RaccoonClient };
