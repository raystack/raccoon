const RaccoonClient = require('../lib/rest');
const axios = require('axios');

const mockHTTPClient = {
    post: jest.fn()
};

jest.mock('uuid', () => ({
    v4: jest.fn().mockReturnValue('mocked-uuid')
}));

describe('RaccoonClient', () => {

    describe('constructor', () => {
        it('should create an instance with default configuration', () => {
            const raccoonClient = new RaccoonClient();

            expect(raccoonClient).toBeDefined();
            expect(raccoonClient.url).toBe('');
            expect(raccoonClient.headers).toStrictEqual({});
            expect(raccoonClient.httpClient).toBeDefined;
            expect(raccoonClient.serializer).toBeDefined;
            expect(raccoonClient.retryMax).toBe(3);
            expect(raccoonClient.retryWait).toBe(5000);
            expect(raccoonClient.wire).toStrictEqual({ ContentType: 'application/json' });
        });

        it('should create an instance with provided configuration', () => {
            const config = {
                serializationType: 'protobuf',
                wire: { ContentType: 'application/custom' },
                headers: { 'Authorization': 'Bearer token' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };

            const raccoonClient = new RaccoonClient(config, mockHTTPClient);

            expect(raccoonClient).toBeDefined();
            expect(raccoonClient.serializer).toBeDefined;
            expect(raccoonClient.wire).toEqual(config.wire);
            expect(raccoonClient.headers).toEqual(config.headers);
            expect(raccoonClient.retryMax).toBe(config.retryMax);
            expect(raccoonClient.retryWait).toBe(config.retryWait);
            expect(raccoonClient.url).toBe(config.url);
            expect(raccoonClient.httpClient).toBe(mockHTTPClient);
        });
    });
    describe('send', () => {

        afterEach(() => {
            jest.spyOn(Date, 'now').mockRestore();
            mockHTTPClient.post.mockReset();
        });

        it('should send single event along with custom headers and return response', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient);

            const now = new Date('2023-08-07T00:00:00Z');
            Date.now = jest.spyOn(Date, 'now').mockReturnValue(now);

            const sentTime = 1691395789;

            const events = [
                { type: 'topic', data: { key1: 'value' } }
            ];

            mockHTTPClient.post.mockResolvedValue({ data: { status: 1, code: 1, sent_time: sentTime, data: { req_guid: 'mocked-uuid' } } });

            const response = await raccoonClient.send(events);

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {
                    Status: 1,
                    Code: 1,
                    SentTime: sentTime,
                    Data: { req_guid: 'mocked-uuid' },
                },
                error: null,
            });

            expect(mockHTTPClient.post).toHaveBeenCalledWith(
                'http://example.com/api',
                JSON.stringify({
                    req_guid: 'mocked-uuid',
                    events: [
                        { type: 'topic', event_bytes: [123, 34, 107, 101, 121, 49, 34, 58, 34, 118, 97, 108, 117, 101, 34, 125] }
                    ],
                    sent_time: { seconds: 1691366400, nanos: 0 },
                }),
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'X-User-ID': 'test-user-1',
                    }
                }
            );
        });

        it('should send multiple events along with custom headers and return response', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient);

            const now = new Date('2023-08-07T00:00:00Z');
            Date.now = jest.spyOn(Date, 'now').mockReturnValue(now);

            const sentTime = 1691395789;

            const events = [
                { type: 'topic1', data: { key1: 'value1' } },
                { type: 'topic2', data: { key2: 'value2' } },
            ];

            mockHTTPClient.post.mockResolvedValue({ data: { status: 1, code: 1, sent_time: sentTime, data: { req_guid: 'mocked-uuid' } } });

            const response = await raccoonClient.send(events);

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {
                    Status: 1,
                    Code: 1,
                    SentTime: sentTime,
                    Data: { req_guid: 'mocked-uuid' },
                },
                error: null,
            });

            expect(mockHTTPClient.post).toHaveBeenCalledWith(
                'http://example.com/api',
                JSON.stringify({
                    req_guid: 'mocked-uuid',
                    events: [
                        { type: 'topic1', event_bytes: [123, 34, 107, 101, 121, 49, 34, 58, 34, 118, 97, 108, 117, 101, 49, 34, 125] },
                        { type: 'topic2', event_bytes: [123, 34, 107, 101, 121, 50, 34, 58, 34, 118, 97, 108, 117, 101, 50, 34, 125] },
                    ],
                    sent_time: { seconds: 1691366400, nanos: 0 },
                }),
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'X-User-ID': 'test-user-1',
                    }
                }
            );
        });

        it('should return error if input message do not have data', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient);

            const invalidEvents = [
                { type: 'topic', daa: { key1: 'value' } }
            ];

            const response = await raccoonClient.send(invalidEvents);

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {},
                error: expect.any(Error),
            });
            expect(response.error.message).toEqual("Invalid event: {\"type\":\"topic\",\"daa\":{\"key1\":\"value\"}}");
            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should return error if input message do not have type', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient);

            const invalidEvents = [
                { tye: 'topic', data: { key1: 'value' } }
            ];

            const response = await raccoonClient.send(invalidEvents);

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {},
                error: expect.any(Error),
            });
            expect(response.error.message).toEqual("Invalid event: {\"tye\":\"topic\",\"data\":{\"key1\":\"value\"}}");
            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should return error if input messages empty', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient);

            const response = await raccoonClient.send([]);

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {},
                error: expect.any(Error),
            });
            expect(response.error.message).toEqual("No events provided");
            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should return error if input messages not sent', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient);

            const response = await raccoonClient.send();

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {},
                error: expect.any(Error),
            });
            expect(response.error.message).toEqual("No events provided");
            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should return error if input messages null', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient);

            const response = await raccoonClient.send(null);

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {},
                error: expect.any(Error),
            });
            expect(response.error.message).toEqual("No events provided");
            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should return error if error from server', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 5,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient);

            const now = new Date('2023-08-07T00:00:00Z');
            Date.now = jest.spyOn(Date, 'now').mockReturnValue(now);

            const sentTime = 1691395789;

            const events = [
                { type: 'topic', data: { key1: 'value' } }
            ];

            mockHTTPClient.post.mockRejectedValue({ status: 400, statusText: 'Bad Request' });

            const response = await raccoonClient.send(events);

            expect(mockHTTPClient.post).toHaveBeenCalledWith(
                'http://example.com/api',
                JSON.stringify({
                    req_guid: 'mocked-uuid',
                    events: [
                        { type: 'topic', event_bytes: [123, 34, 107, 101, 121, 49, 34, 58, 34, 118, 97, 108, 117, 101, 34, 125] }
                    ],
                    sent_time: { seconds: 1691366400, nanos: 0 },
                }),
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'X-User-ID': 'test-user-1',
                    }
                }
            );

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {},
                error: {
                    status: 400,
                    statusText: 'Bad Request',
                },
            });
        });
    });
});
