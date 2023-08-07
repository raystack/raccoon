const axios = require('axios');
const uuid = require('uuid');

const createJsonSerializer = require('./serializer/json_serializer');
const createProtobufSerializer = require('./serializer/proto_serializer');
const retry = require('./retry/retry')

const NANOSECONDS_PER_MILLISECOND = 1e6;

class RaccoonClient {
    /**
     * Creates a new instance of the RaccoonClient.
     *
     * @constructor
     * @param {Object} config - Configuration options for the RaccoonClient.
     * @param {string} [config.serializationType='json'] - The serialization type to use, either 'protobuf' or 'json'.
     * @param {Object} [config.wire] - The wire configuration, containing options like ContentType.
     * @param {Object} [config.headers] - Custom headers to be included in the HTTP requests.
     * @param {number} [config.retryMax=3] - The maximum number of retry attempts for failed requests.
     * @param {number} [config.retryWait=1000] - The time in milliseconds to wait between retry attempts.
     * @param {string} [config.url=''] - The base URL for the API requests.
     * @returns {RaccoonClient} A new instance of the RaccoonClient.
     */
    constructor(config = {}, httpClient) {
        this.serializer = config.serializationType === 'protobuf'
            ? createProtobufSerializer()
            : createJsonSerializer();
        this.wire = config.wire || { ContentType: 'application/json' };
        this.httpClient = httpClient || axios.create();
        this.headers = config.headers || {};
        this.retryMax = config.retryMax || 3;
        this.retryWait = config.retryWait || 5000;
        this.url = config.url || '';
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
        let requestId;
        try {
            requestId = uuid.v4();
            if (!events || events.length === 0) {
                throw new Error('No events provided');
            }
            console.info(`started request, url: ${this.url}, req-id: ${requestId}`);
            const eventsToSend = [];
            for (const event of events) {
                if (event && event.type && event.data) {
                    eventsToSend.push({
                        type: event.type,
                        event_bytes: this.serializer.serialize(event.data),
                    });
                } else {
                    throw new Error(`Invalid event: ${JSON.stringify(event)}`);
                }
            }

            const sentTime = this.getCurrentTimestamp();

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

            console.info(`ended request, url: ${this.url}, req-id: ${requestId}`);
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
            console.error(`error, url: ${this.url}, req-id: ${requestId}, ${error}`);
            return { reqId: requestId, response: {}, error };
        }
    }

    async executeRequest(raccoonRequest) {
        this.headers['Content-Type'] = this.wire.ContentType;
        const response = await this.httpClient.post(this.url, raccoonRequest, {
            headers: this.headers,
        });
        return response.data;
    }

    getCurrentTimestamp() {
        const now = Date.now();
        return ({
            seconds: Math.floor(now / 1000),
            nanos: (now % 1000) * NANOSECONDS_PER_MILLISECOND
        });
    }
}

module.exports = RaccoonClient;
