import { RaccoonClient } from '../lib/rest.js';
import { jest } from '@jest/globals'

const mockHTTPClient = {
    post: jest.fn()
};

const mockUUIDGenerator = (() => 'mocked-uuid');

describe('RaccoonClient', () => {

    describe('constructor', () => {
        it('should create an instance with default configuration', () => {
            const raccoonClient = new RaccoonClient();

            expect(raccoonClient).toBeDefined();
            expect(raccoonClient.url).toBe('');
            expect(raccoonClient.headers).toStrictEqual({});
            expect(raccoonClient.httpClient).toBeDefined;
            expect(raccoonClient.serializer).toBeDefined;
            expect(raccoonClient.logger).toBe(console);
            expect(raccoonClient.retryMax).toBe(3);
            expect(raccoonClient.retryWait).toBe(5000);
            expect(raccoonClient.wire).toStrictEqual({ ContentType: 'application/json' });
        });

        it('should create an instance with provided configuration', () => {

            const mockLogger = {
                info: jest.fn(),
                error: jest.fn(),
            };

            const options = {
                serializationType: 'protobuf',
                wire: { ContentType: 'application/custom' },
                headers: { 'Authorization': 'Bearer token' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
                logger: mockLogger,
            };

            const raccoonClient = new RaccoonClient(options, mockHTTPClient);

            expect(raccoonClient).toBeDefined();
            expect(raccoonClient.serializer).toBeDefined;
            expect(raccoonClient.wire).toEqual(options.wire);
            expect(raccoonClient.headers).toEqual(options.headers);
            expect(raccoonClient.retryMax).toBe(options.retryMax);
            expect(raccoonClient.retryWait).toBe(options.retryWait);
            expect(raccoonClient.url).toBe(options.url);
            expect(raccoonClient.httpClient).toBe(mockHTTPClient);
            expect(raccoonClient.logger).toBe(mockLogger)
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
            const raccoonClient = new RaccoonClient(config, mockHTTPClient, mockUUIDGenerator);

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
            const raccoonClient = new RaccoonClient(config, mockHTTPClient, mockUUIDGenerator);

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

        it('should throw error if input message do not have data', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient, mockUUIDGenerator);

            const invalidEvents = [
                { type: 'topic', daa: { key1: 'value' } }
            ];

            await expect(() => raccoonClient.send(invalidEvents)).rejects.toThrow("req-id: mocked-uuid, error: Error: Invalid event: {\"type\":\"topic\",\"daa\":{\"key1\":\"value\"}}");
            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should throw error if input message do not have type', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient, mockUUIDGenerator);

            const invalidEvents = [
                { tye: 'topic', data: { key1: 'value' } }
            ];

            await expect(() => raccoonClient.send(invalidEvents)).rejects.toThrow("req-id: mocked-uuid, error: Error: Invalid event: {\"tye\":\"topic\",\"data\":{\"key1\":\"value\"}}");
            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should throw error if input messages empty', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient, mockUUIDGenerator);

            await expect(() => raccoonClient.send([])).rejects.toThrow("req-id: mocked-uuid, error: Error: No events provided");

            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should throw error if input messages not sent', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient, mockUUIDGenerator);

            await expect(() => raccoonClient.send()).rejects.toThrow("req-id: mocked-uuid, error: Error: No events provided");

            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should throw error if input messages null', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient, mockUUIDGenerator);

            await expect(() => raccoonClient.send(null)).rejects.toThrow("req-id: mocked-uuid, error: Error: No events provided");
            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should throw error if error from server', async () => {
            const config = {
                serializationType: 'json',
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 5,
                url: 'http://example.com/api',
            };
            const raccoonClient = new RaccoonClient(config, mockHTTPClient, mockUUIDGenerator);

            const now = new Date('2023-08-07T00:00:00Z');
            Date.now = jest.spyOn(Date, 'now').mockReturnValue(now);

            const events = [
                { type: 'topic', data: { key1: 'value' } }
            ];

            mockHTTPClient.post.mockRejectedValue(JSON.stringify({ status: 400, statusText: 'Bad Request' }));

            await expect(() => raccoonClient.send(events)).rejects.toThrow("req-id: mocked-uuid, error: {\"status\":400,\"statusText\":\"Bad Request\"}");

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
    });
});
