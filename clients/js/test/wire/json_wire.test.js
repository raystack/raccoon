import { createJsonMarshaller } from '../../lib/wire/json_wire.js';
import { WireType } from '../../lib/rest.js';
import protos from '../../protos/proton_compiled.js';

describe('JsonMarshaller', () => {
    const jsonMarshaller = createJsonMarshaller();

    describe('marshal', () => {
        it('should marshal a valid data object', () => {
            const Event = protos.raystack.raccoon.v1beta1.Event;
            const eventMessage = new Event();
            eventMessage.type = 'test-type';
            eventMessage.event_bytes = Buffer.from('test-data');

            const result = jsonMarshaller.marshal(eventMessage);

            expect(result).toEqual(
                '{"event_bytes":{"type":"Buffer","data":[116,101,115,116,45,100,97,116,97]},"type":"test-type"}'
            );
        });

        it('should throw an error for invalid data object', () => {
            expect(() => jsonMarshaller.marshal(null)).toThrow(
                'Invalid data object for marshalling'
            );
            expect(() => jsonMarshaller.marshal(undefined)).toThrow(
                'Invalid data object for marshalling'
            );
        });

        it('should throw an error if data object constructor is missing toObject method', () => {
            const mockData = {};

            expect(() => jsonMarshaller.marshal(mockData)).toThrow(
                'Invalid Protobuf message object for marshalling'
            );
        });

        it('should throw an error if data object constructor toObject is not a function', () => {
            const mockData = {
                constructor: {
                    toObject: 'not_a_function'
                }
            };

            expect(() => jsonMarshaller.marshal(mockData)).toThrow(
                'Invalid Protobuf message object for marshalling'
            );
        });
    });

    describe('unmarshal', () => {
        it('should unmarshal data using a valid type', () => {
            const Event = protos.raystack.raccoon.v1beta1.Event;
            const eventObject = {
                type: 'test-type',
                event_bytes: Buffer.from('test-data')
            };

            const result = jsonMarshaller.unmarshal(eventObject, Event);

            expect(result.type).toEqual('test-type');
            expect(result.event_bytes).toEqual(Buffer.from('test-data'));
        });

        it('should throw an error for invalid data object', () => {
            expect(() => jsonMarshaller.unmarshal(null, {})).toThrow(
                'Invalid data object for unmarshalling'
            );
            expect(() => jsonMarshaller.unmarshal(undefined, {})).toThrow(
                'Invalid data object for unmarshalling'
            );
        });

        it('should throw an error for null type', () => {
            const mockData = { field1: 'value1', field2: 'value2' };

            expect(() => jsonMarshaller.unmarshal(mockData, null)).toThrow(
                'Invalid Protobuf message type for unmarshalling'
            );
        });

        it('should throw an error for undefined type', () => {
            const mockData = { field1: 'value1', field2: 'value2' };

            expect(() => jsonMarshaller.unmarshal(mockData, undefined)).toThrow(
                'Invalid Protobuf message type for unmarshalling'
            );
        });

        it('should throw an error if type.fromObject is not a function', () => {
            const mockData = { field1: 'value1', field2: 'value2' };
            const invalidType = {
                fromObject: 'not_a_function'
            };

            expect(() => jsonMarshaller.unmarshal(mockData, invalidType)).toThrow(
                'Invalid Protobuf message type for unmarshalling'
            );
        });
    });

    describe('getContentType', () => {
        it('should return WireType.JSON', () => {
            const result = jsonMarshaller.getContentType();
            expect(result).toEqual(WireType.JSON);
        });
    });

    describe('getResponseType', () => {
        it('should return "json"', () => {
            const result = jsonMarshaller.getResponseType();
            expect(result).toEqual('json');
        });
    });
});
