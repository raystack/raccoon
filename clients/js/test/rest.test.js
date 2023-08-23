// eslint-disable-next-line
import { jest } from '@jest/globals';
import { RaccoonClient, SerializationType, WireType } from '../lib/rest.js';
import { raystack, google } from '../protos/proton_compiled.js';

const mockHTTPClient = {
    post: jest.fn()
};

const mockUUIDGenerator = () => 'mocked-uuid';

describe('RaccoonClient', () => {
    describe('constructor', () => {
        it('should create an instance with provided configuration', () => {
            const mockLogger = {
                info: jest.fn(),
                error: jest.fn()
            };

            const options = {
                serializationType: SerializationType.JSON,
                wireType: WireType.JSON,
                headers: { Authorization: 'Bearer token' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
                logger: mockLogger,
                timeout: 1000
            };

            const raccoonClient = new RaccoonClient(options);

            expect(raccoonClient).toBeDefined();
            expect(raccoonClient.serialize).toBeDefined();
            expect(raccoonClient.marshaller).toBeDefined();
            expect(raccoonClient.headers).toEqual(options.headers);
            expect(raccoonClient.retryMax).toBe(options.retryMax);
            expect(raccoonClient.retryWait).toBe(options.retryWait);
            expect(raccoonClient.url).toBe(options.url);
            expect(raccoonClient.logger).toBe(mockLogger);
            expect(raccoonClient.timeout).toBe(options.timeout);
        });

        it('should add default values', () => {
            const options = {
                serializationType: SerializationType.JSON,
                wireType: WireType.JSON
            };

            const raccoonClient = new RaccoonClient(options);

            expect(raccoonClient).toBeDefined();
            expect(raccoonClient.serialize).toBeDefined();
            expect(raccoonClient.marshaller).toBeDefined();
            expect(raccoonClient.headers).toEqual({});
            expect(raccoonClient.retryMax).toBe(3);
            expect(raccoonClient.retryWait).toBe(5000);
            expect(raccoonClient.logger).toBe(console);
            expect(raccoonClient.timeout).toBe(5000);
        });

        it('should throw error for invalid serializationType', () => {
            const mockLogger = {
                info: jest.fn(),
                error: jest.fn()
            };

            const options = {
                serializationType: 'invalidType',
                wireType: WireType.JSON,
                headers: { Authorization: 'Bearer token' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
                logger: mockLogger,
                timeout: 1000
            };

            expect(() => new RaccoonClient(options)).toThrow(
                'Invalid serializationType: invalidType'
            );
        });

        it('should throw error for invalid wireType', () => {
            const mockLogger = {
                info: jest.fn(),
                error: jest.fn()
            };

            const options = {
                serializationType: SerializationType.PROTOBUF,
                wireType: 'invalidType',
                headers: { Authorization: 'Bearer token' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api',
                logger: mockLogger,
                timeout: 1000
            };

            expect(() => new RaccoonClient(options)).toThrow('Invalid wireType: invalidType');
        });
    });
    describe('send', () => {
        afterEach(() => {
            jest.spyOn(Date, 'now').mockRestore();
            mockHTTPClient.post.mockReset();
        });

        it('should send JSON event along with custom headers and return response with JSON wire type', async () => {
            const config = {
                serializationType: SerializationType.JSON,
                wireType: WireType.JSON,
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api'
            };
            const raccoonClient = new RaccoonClient(config);
            raccoonClient.httpClient = mockHTTPClient;
            raccoonClient.uuidGenerator = mockUUIDGenerator;

            const now = new Date('2023-08-07T00:00:00Z');
            Date.now = jest.spyOn(Date, 'now').mockReturnValue(now);

            const sentTime = 1691366400;

            const events = [{ type: 'topic', data: { key1: 'value' } }];

            mockHTTPClient.post.mockResolvedValue({
                data: {
                    status: 1,
                    code: 1,
                    sent_time: sentTime,
                    data: { req_guid: 'mocked-uuid' }
                }
            });

            const response = await raccoonClient.send(events);

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {
                    status: 'STATUS_SUCCESS',
                    code: 'CODE_OK',
                    sent_time: '1691366400',
                    data: { req_guid: 'mocked-uuid' }
                },
                error: null
            });

            expect(mockHTTPClient.post).toHaveBeenCalledWith(
                'http://example.com/api',
                JSON.stringify({
                    req_guid: 'mocked-uuid',
                    sent_time: { seconds: 1691366400, nanos: 0 },
                    events: [
                        {
                            event_bytes: [
                                123, 34, 107, 101, 121, 49, 34, 58, 34, 118, 97, 108, 117, 101, 34,
                                125
                            ],
                            type: 'topic'
                        }
                    ]
                }),
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'X-User-ID': 'test-user-1'
                    },
                    timeout: 5000,
                    responseType: 'json'
                }
            );
        });

        it('should send JSON event along with custom headers and return response with Protobuf wire type', async () => {
            const config = {
                serializationType: SerializationType.JSON,
                wireType: WireType.PROTOBUF,
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api'
            };
            const raccoonClient = new RaccoonClient(config);
            raccoonClient.httpClient = mockHTTPClient;
            raccoonClient.uuidGenerator = mockUUIDGenerator;

            const now = new Date('2023-08-07T00:00:00Z');
            Date.now = jest.spyOn(Date, 'now').mockReturnValue(now);

            const sentTime = 1691366400;

            const events = [{ type: 'topic', data: { key1: 'value' } }];

            const { SendEventRequest } = raystack.raccoon.v1beta1;
            const { SendEventResponse } = raystack.raccoon.v1beta1;

            mockHTTPClient.post.mockResolvedValue({
                data: SendEventResponse.encode(
                    SendEventResponse.create({
                        status: 1,
                        code: 1,
                        sent_time: sentTime,
                        data: { req_guid: 'mocked-uuid' }
                    })
                ).finish()
            });

            const response = await raccoonClient.send(events);

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {
                    status: 'STATUS_SUCCESS',
                    code: 'CODE_OK',
                    sent_time: '1691366400',
                    data: { req_guid: 'mocked-uuid' }
                },
                error: null
            });

            expect(mockHTTPClient.post).toHaveBeenCalledWith(
                'http://example.com/api',
                new Uint8Array(
                    SendEventRequest.encode(
                        SendEventRequest.create({
                            req_guid: 'mocked-uuid',
                            sent_time: { seconds: 1691366400, nanos: 0 },
                            events: [
                                {
                                    event_bytes: [
                                        123, 34, 107, 101, 121, 49, 34, 58, 34, 118, 97, 108, 117,
                                        101, 34, 125
                                    ],
                                    type: 'topic'
                                }
                            ]
                        })
                    ).finish()
                ),
                {
                    headers: {
                        'Content-Type': 'application/proto',
                        'X-User-ID': 'test-user-1'
                    },
                    timeout: 5000,
                    responseType: 'arraybuffer'
                }
            );
        });

        it('should send Protobuf event along with custom headers and return response with JSON wire type', async () => {
            const config = {
                serializationType: SerializationType.PROTOBUF,
                wireType: WireType.JSON,
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api'
            };
            const raccoonClient = new RaccoonClient(config);
            raccoonClient.httpClient = mockHTTPClient;
            raccoonClient.uuidGenerator = mockUUIDGenerator;

            const now = new Date('2023-08-07T00:00:00Z');
            Date.now = jest.spyOn(Date, 'now').mockReturnValue(now);

            const sentTime = 1691366400;

            const timestamp = new google.protobuf.Timestamp();
            timestamp.seconds = 123;
            timestamp.nanos = 12;

            const events = [
                {
                    type: 'topic',
                    data: timestamp
                }
            ];

            mockHTTPClient.post.mockResolvedValue({
                data: {
                    status: 1,
                    code: 1,
                    sent_time: sentTime,
                    data: { req_guid: 'mocked-uuid' }
                }
            });

            const response = await raccoonClient.send(events);

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {
                    status: 'STATUS_SUCCESS',
                    code: 'CODE_OK',
                    sent_time: '1691366400',
                    data: { req_guid: 'mocked-uuid' }
                },
                error: null
            });

            expect(mockHTTPClient.post).toHaveBeenCalledWith(
                'http://example.com/api',
                JSON.stringify({
                    req_guid: 'mocked-uuid',
                    sent_time: { seconds: 1691366400, nanos: 0 },
                    events: [
                        {
                            event_bytes: [8, 123, 16, 12],
                            type: 'topic'
                        }
                    ]
                }),
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'X-User-ID': 'test-user-1'
                    },
                    timeout: 5000,
                    responseType: 'json'
                }
            );
        });

        it('should send Protobuf event along with custom headers and return response with Protobuf wire type', async () => {
            const config = {
                serializationType: SerializationType.PROTOBUF,
                wireType: WireType.PROTOBUF,
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api'
            };
            const raccoonClient = new RaccoonClient(config);
            raccoonClient.httpClient = mockHTTPClient;
            raccoonClient.uuidGenerator = mockUUIDGenerator;

            const now = new Date('2023-08-07T00:00:00Z');
            Date.now = jest.spyOn(Date, 'now').mockReturnValue(now);

            const sentTime = 1691366400;

            const timestamp = new google.protobuf.Timestamp();
            timestamp.seconds = 123;
            timestamp.nanos = 12;

            const events = [
                {
                    type: 'topic',
                    data: timestamp
                }
            ];

            const { SendEventRequest } = raystack.raccoon.v1beta1;
            const { SendEventResponse } = raystack.raccoon.v1beta1;

            mockHTTPClient.post.mockResolvedValue({
                data: SendEventResponse.encode(
                    SendEventResponse.create({
                        status: 1,
                        code: 1,
                        sent_time: sentTime,
                        data: { req_guid: 'mocked-uuid' }
                    })
                ).finish()
            });

            const response = await raccoonClient.send(events);

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {
                    status: 'STATUS_SUCCESS',
                    code: 'CODE_OK',
                    sent_time: '1691366400',
                    data: { req_guid: 'mocked-uuid' }
                },
                error: null
            });

            expect(mockHTTPClient.post).toHaveBeenCalledWith(
                'http://example.com/api',
                new Uint8Array(
                    SendEventRequest.encode(
                        SendEventRequest.create({
                            req_guid: 'mocked-uuid',
                            sent_time: { seconds: 1691366400, nanos: 0 },
                            events: [
                                {
                                    event_bytes: [8, 123, 16, 12],
                                    type: 'topic'
                                }
                            ]
                        })
                    ).finish()
                ),
                {
                    headers: {
                        'Content-Type': 'application/proto',
                        'X-User-ID': 'test-user-1'
                    },
                    timeout: 5000,
                    responseType: 'arraybuffer'
                }
            );
        });

        it('should send multiple events along with custom headers and return response', async () => {
            const config = {
                serializationType: SerializationType.JSON,
                wireType: WireType.JSON,
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api'
            };
            const raccoonClient = new RaccoonClient(config);
            raccoonClient.httpClient = mockHTTPClient;
            raccoonClient.uuidGenerator = mockUUIDGenerator;

            const now = new Date('2023-08-07T00:00:00Z');
            Date.now = jest.spyOn(Date, 'now').mockReturnValue(now);

            const sentTime = 1691366400;

            const events = [
                { type: 'topic1', data: { key1: 'value1' } },
                { type: 'topic2', data: { key2: 'value2' } }
            ];

            mockHTTPClient.post.mockResolvedValue({
                data: {
                    status: 1,
                    code: 1,
                    sent_time: sentTime,
                    data: { req_guid: 'mocked-uuid' }
                }
            });

            const response = await raccoonClient.send(events);

            expect(response).toEqual({
                reqId: 'mocked-uuid',
                response: {
                    status: 'STATUS_SUCCESS',
                    code: 'CODE_OK',
                    sent_time: '1691366400',
                    data: { req_guid: 'mocked-uuid' }
                },
                error: null
            });

            expect(mockHTTPClient.post).toHaveBeenCalledWith(
                'http://example.com/api',
                JSON.stringify({
                    req_guid: 'mocked-uuid',
                    sent_time: { seconds: 1691366400, nanos: 0 },
                    events: [
                        {
                            event_bytes: [
                                123, 34, 107, 101, 121, 49, 34, 58, 34, 118, 97, 108, 117, 101, 49,
                                34, 125
                            ],
                            type: 'topic1'
                        },
                        {
                            event_bytes: [
                                123, 34, 107, 101, 121, 50, 34, 58, 34, 118, 97, 108, 117, 101, 50,
                                34, 125
                            ],
                            type: 'topic2'
                        }
                    ]
                }),
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'X-User-ID': 'test-user-1'
                    },
                    timeout: 5000,
                    responseType: 'json'
                }
            );
        });

        it('should throw error if input message do not have data', async () => {
            const config = {
                serializationType: SerializationType.JSON,
                wireType: WireType.JSON,
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api'
            };
            const raccoonClient = new RaccoonClient(config);
            raccoonClient.httpClient = mockHTTPClient;
            raccoonClient.uuidGenerator = mockUUIDGenerator;

            const invalidEvents = [{ type: 'topic', daa: { key1: 'value' } }];

            await expect(() => raccoonClient.send(invalidEvents)).rejects.toThrow(
                'req-id: mocked-uuid, error: Error: Invalid event: {"type":"topic","daa":{"key1":"value"}}'
            );
            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should throw error if input message do not have type', async () => {
            const config = {
                serializationType: SerializationType.JSON,
                wireType: WireType.JSON,
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api'
            };
            const raccoonClient = new RaccoonClient(config);
            raccoonClient.httpClient = mockHTTPClient;
            raccoonClient.uuidGenerator = mockUUIDGenerator;

            const invalidEvents = [{ tye: 'topic', data: { key1: 'value' } }];

            await expect(() => raccoonClient.send(invalidEvents)).rejects.toThrow(
                'req-id: mocked-uuid, error: Error: Invalid event: {"tye":"topic","data":{"key1":"value"}}'
            );
            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should throw error if input messages empty', async () => {
            const config = {
                serializationType: SerializationType.JSON,
                wireType: WireType.JSON,
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api'
            };
            const raccoonClient = new RaccoonClient(config);
            raccoonClient.httpClient = mockHTTPClient;
            raccoonClient.uuidGenerator = mockUUIDGenerator;

            await expect(() => raccoonClient.send([])).rejects.toThrow(
                'req-id: mocked-uuid, error: Error: No events provided'
            );

            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should throw error if input messages not sent', async () => {
            const config = {
                serializationType: SerializationType.JSON,
                wireType: WireType.JSON,
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api'
            };
            const raccoonClient = new RaccoonClient(config);
            raccoonClient.httpClient = mockHTTPClient;
            raccoonClient.uuidGenerator = mockUUIDGenerator;

            await expect(() => raccoonClient.send()).rejects.toThrow(
                'req-id: mocked-uuid, error: Error: No events provided'
            );

            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should throw error if input messages null', async () => {
            const config = {
                serializationType: SerializationType.JSON,
                wireType: WireType.JSON,
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 3000,
                url: 'http://example.com/api'
            };
            const raccoonClient = new RaccoonClient(config);
            raccoonClient.httpClient = mockHTTPClient;
            raccoonClient.uuidGenerator = mockUUIDGenerator;

            await expect(() => raccoonClient.send(null)).rejects.toThrow(
                'req-id: mocked-uuid, error: Error: No events provided'
            );
            expect(mockHTTPClient.post).not.toHaveBeenCalled();
        });

        it('should throw error if error from server', async () => {
            const config = {
                serializationType: SerializationType.JSON,
                wireType: WireType.JSON,
                headers: { 'X-User-ID': 'test-user-1' },
                retryMax: 5,
                retryWait: 5,
                url: 'http://example.com/api'
            };
            const raccoonClient = new RaccoonClient(config);
            raccoonClient.httpClient = mockHTTPClient;
            raccoonClient.uuidGenerator = mockUUIDGenerator;

            const now = new Date('2023-08-07T00:00:00Z');
            Date.now = jest.spyOn(Date, 'now').mockReturnValue(now);

            const events = [{ type: 'topic', data: { key1: 'value' } }];

            mockHTTPClient.post.mockRejectedValue(
                JSON.stringify({ status: 400, statusText: 'Bad Request' })
            );

            await expect(() => raccoonClient.send(events)).rejects.toThrow(
                'req-id: mocked-uuid, error: {"status":400,"statusText":"Bad Request"}'
            );

            expect(mockHTTPClient.post).toHaveBeenCalledWith(
                'http://example.com/api',
                JSON.stringify({
                    req_guid: 'mocked-uuid',
                    sent_time: { seconds: 1691366400, nanos: 0 },
                    events: [
                        {
                            event_bytes: [
                                123, 34, 107, 101, 121, 49, 34, 58, 34, 118, 97, 108, 117, 101, 34,
                                125
                            ],
                            type: 'topic'
                        }
                    ]
                }),
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'X-User-ID': 'test-user-1'
                    },
                    timeout: 5000,
                    responseType: 'json'
                }
            );
        });
    });
});
