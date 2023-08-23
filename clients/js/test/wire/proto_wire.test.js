import { createProtoMarshaller } from '../../lib/wire/proto_wire.js';
import { WireType } from '../../lib/rest.js';
import protos from '../../protos/proton_compiled.js';

describe('ProtoMarshaller', () => {
    const protoMarshaller = createProtoMarshaller();

    describe('marshal', () => {
        it('should marshal a valid data object', () => {
            const Event = protos.raystack.raccoon.v1beta1.Event;
            const eventMessage = new Event();
            eventMessage.type = 'test-type';
            eventMessage.event_bytes = Buffer.from('test-data');

            const result = protoMarshaller.marshal(eventMessage);
            const decodedData = Event.decode(result);

            expect(decodedData.type).toEqual('test-type');
            expect(decodedData.event_bytes).toEqual(
                new Uint8Array([116, 101, 115, 116, 45, 100, 97, 116, 97])
            );
        });

        it('should throw an error for invalid data object', () => {
            expect(() => protoMarshaller.marshal(null)).toThrow(
                'Invalid data object for marshalling'
            );
            expect(() => protoMarshaller.marshal(undefined)).toThrow(
                'Invalid data object for marshalling'
            );
        });

        it('should throw an error if data object constructor is missing encode method', () => {
            const mockData = {};

            expect(() => protoMarshaller.marshal(mockData)).toThrow(
                'Invalid Protobuf message object for marshalling'
            );
        });

        it('should throw an error if data object constructor encode is not a function', () => {
            const mockData = {
                constructor: {
                    encode: 'not_a_function'
                }
            };

            expect(() => protoMarshaller.marshal(mockData)).toThrow(
                'Invalid Protobuf message object for marshalling'
            );
        });
    });

    describe('unmarshal', () => {
        it('should unmarshal data using a valid type', () => {
            const Event = protos.raystack.raccoon.v1beta1.Event;
            const eventMessage = new Event();
            eventMessage.type = 'test-type';
            eventMessage.event_bytes = Buffer.from('test-data');

            const bytes = Event.encode(eventMessage).finish();

            const result = protoMarshaller.unmarshal(bytes, Event);

            expect(result.type).toEqual('test-type');
            expect(result.event_bytes).toEqual(
                new Uint8Array([116, 101, 115, 116, 45, 100, 97, 116, 97])
            );
        });

        it('should throw an error for invalid data', () => {
            expect(() => protoMarshaller.unmarshal(null, {})).toThrow(
                'Invalid data for unmarshalling'
            );
            expect(() => protoMarshaller.unmarshal(undefined, {})).toThrow(
                'Invalid data for unmarshalling'
            );
        });

        it('should throw an error for invalid type', () => {
            const mockData = new Uint8Array([1, 2, 3]);

            expect(() => protoMarshaller.unmarshal(mockData, null)).toThrow(
                'Invalid Protobuf message type for unmarshalling'
            );
            expect(() => protoMarshaller.unmarshal(mockData, undefined)).toThrow(
                'Invalid Protobuf message type for unmarshalling'
            );
        });

        it('should throw an error if type.decode is not a function', () => {
            const mockData = new Uint8Array([1, 2, 3]);
            const invalidType = {
                decode: 'not_a_function'
            };

            expect(() => protoMarshaller.unmarshal(mockData, invalidType)).toThrow(
                'Invalid Protobuf message type for unmarshalling'
            );
        });
    });

    describe('getContentType', () => {
        it('should return WireType.PROTOBUF', () => {
            const result = protoMarshaller.getContentType();
            expect(result).toEqual(WireType.PROTOBUF);
        });
    });

    describe('getResponseType', () => {
        it('should return "arraybuffer"', () => {
            const result = protoMarshaller.getResponseType();
            expect(result).toEqual('arraybuffer');
        });
    });
});
