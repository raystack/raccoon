import { createJsonSerializer } from '../../lib/serializer/json_serializer.js';

describe('JsonSerializer', () => {
    test('should convert data to JSON', () => {
        const serialize = createJsonSerializer();
        const data = { key: 'value' };
        const serializedData = serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle empty object', () => {
        const serialize = createJsonSerializer();
        const data = {};
        const serializedData = serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle empty array', () => {
        const serialize = createJsonSerializer();
        const data = [];
        const serializedData = serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle null', () => {
        const serialize = createJsonSerializer();
        const data = null;

        expect(() => {
            serialize(data);
        }).toThrowError("Invalid data object for serialization");
    });

    test('should handle undefined with proper error', () => {
        const serialize = createJsonSerializer();
        const data = undefined;

        expect(() => {
            serialize(data);
        }).toThrowError("Invalid data object for serialization");
    });


    test('should handle numbers', () => {
        const serialize = createJsonSerializer();
        const data = 12345;
        const serializedData = serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle strings', () => {
        const serialize = createJsonSerializer();
        const data = 'hello';
        const serializedData = serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });
});
