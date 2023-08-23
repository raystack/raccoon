import { createProtobufSerializer } from '../../lib/serializer/proto_serializer.js';
import protos from '../../protos/proton_compiled.js';

describe('ProtoSerializer', () => {
    test('should convert data to Proto object', () => {
        const serialize = createProtobufSerializer();
        const Event = protos.raystack.raccoon.v1beta1.Event;
        const eventMessage = new Event();
        eventMessage.type = 'test-type';
        eventMessage.event_bytes = Buffer.from('test-data');

        const serializedData = serialize(eventMessage);
        const decodedData = Event.decode(serializedData);

        expect(decodedData.type).toEqual('test-type');
        expect(decodedData.event_bytes).toEqual([116, 101, 115, 116, 45, 100, 97, 116, 97]);
    });

    test('should throw error for empty object', () => {
        const serialize = createProtobufSerializer();
        const data = {};

        expect(() => {
            serialize(data);
        }).toThrowError('Invalid Protobuf message object for serialization');
    });

    test('should throw error for null', () => {
        const serialize = createProtobufSerializer();
        const data = null;

        expect(() => {
            serialize(data);
        }).toThrowError('Invalid data object for serialization');
    });

    test('should handle undefined with proper error', () => {
        const serialize = createProtobufSerializer();
        const data = undefined;

        expect(() => {
            serialize(data);
        }).toThrowError('Invalid data object for serialization');
    });

    it('should throw an error if data object constructor is missing encode method', () => {
        const mockData = {
            constructor: {}
        };
        const serialize = createProtobufSerializer();

        expect(() => serialize(mockData)).toThrow(
            'Invalid Protobuf message object for serialization'
        );
    });

    it('should throw an error if data object constructor encode is not a function', () => {
        const mockData = {
            constructor: {
                encode: 'not_a_function'
            }
        };
        const serialize = createProtobufSerializer();

        expect(() => serialize(mockData)).toThrow(
            'Invalid Protobuf message object for serialization'
        );
    });
});
